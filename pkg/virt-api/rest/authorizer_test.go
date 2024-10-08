/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Red Hat, Inc.
 *
 */

package rest

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"

	restful "github.com/emicklei/go-restful"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	authorizationclient "k8s.io/client-go/kubernetes/typed/authorization/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("Authorizer", func() {

	Describe("VirtualMachineInstance Subresources", func() {
		var server *ghttp.Server
		var req *restful.Request

		fakecert := &x509.Certificate{}

		app := authorizor{}
		BeforeEach(func() {
			req = &restful.Request{}
			req.Request = &http.Request{}
			req.Request.URL = &url.URL{}
			req.Request.Header = make(map[string][]string)
			req.Request.Header[userHeader] = []string{"user"}
			req.Request.Header[groupHeader] = []string{"userGroup"}
			req.Request.Header[userExtraHeaderPrefix+"test"] = []string{"userExtraValue"}

			server = ghttp.NewServer()
			config, err := clientcmd.BuildConfigFromFlags(server.URL(), "")
			Expect(err).ToNot(HaveOccurred())

			client, err := authorizationclient.NewForConfig(config)
			Expect(err).ToNot(HaveOccurred())

			app.subjectAccessReview = client.SubjectAccessReviews()
			app.userHeaders = append(app.userHeaders, userHeader)
			app.groupHeaders = append(app.groupHeaders, groupHeader)
			app.userExtraHeaderPrefixes = append(app.userExtraHeaderPrefixes, userExtraHeaderPrefix)
		})

		Context("Subresource api with namespaced resource", func() {
			Context("with namespaced resource", func() {
				BeforeEach(func() {
					req.Request.Method = http.MethodGet
					req.Request.URL.Path = "/apis/subresources.kubevirt.io/v1alpha3/namespaces/default/virtualmachineinstances/testvmi/console"
				})

				It("should reject unauthenticated user", func() {
					allowed, reason, err := app.Authorize(req)
					Expect(err).ToNot(HaveOccurred())
					Expect(allowed).To(BeFalse())
					Expect(reason).To(Equal("request is not authenticated"))
				})

				It("should reject unauthorized user", func() {

					req.Request.TLS = &tls.ConnectionState{}
					req.Request.TLS.PeerCertificates = append(req.Request.TLS.PeerCertificates, fakecert)

					result, err := app.generateAccessReview(req)
					Expect(err).ToNot(HaveOccurred())
					result.Status.Allowed = false
					result.Status.Reason = "just because"

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("POST", "/apis/authorization.k8s.io/v1/subjectaccessreviews"),
							ghttp.RespondWithJSONEncoded(http.StatusOK, result),
						),
					)

					allowed, reason, err := app.Authorize(req)
					Expect(err).ToNot(HaveOccurred())
					Expect(allowed).To(BeFalse())
					Expect(reason).To(Equal("just because"))
				})

				It("should allow authorized user", func() {

					req.Request.TLS = &tls.ConnectionState{}
					req.Request.TLS.PeerCertificates = append(req.Request.TLS.PeerCertificates, fakecert)

					result, err := app.generateAccessReview(req)
					Expect(err).ToNot(HaveOccurred())
					result.Status.Allowed = true

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("POST", "/apis/authorization.k8s.io/v1/subjectaccessreviews"),
							ghttp.RespondWithJSONEncoded(http.StatusOK, result),
						),
					)

					allowed, _, err := app.Authorize(req)
					Expect(err).ToNot(HaveOccurred())
					Expect(allowed).To(BeTrue())
				})

				It("should not allow user if auth check fails", func() {

					req.Request.TLS = &tls.ConnectionState{}
					req.Request.TLS.PeerCertificates = append(req.Request.TLS.PeerCertificates, fakecert)

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("POST", "/apis/authorization.k8s.io/v1/subjectaccessreviews"),
							ghttp.RespondWithJSONEncoded(http.StatusInternalServerError, nil),
						),
					)

					allowed, _, err := app.Authorize(req)
					Expect(err).To(HaveOccurred())
					Expect(allowed).To(BeFalse())
				})
			})

			Context("with namespaced base resource", func() {
				BeforeEach(func() {
					req.Request.Method = http.MethodPut
					req.Request.URL.Path = "/apis/subresources.kubevirt.io/v1alpha3/namespaces/default/expand-vm-spec"
				})

				It("should reject unauthenticated user", func() {
					allowed, reason, err := app.Authorize(req)
					Expect(err).ToNot(HaveOccurred())
					Expect(allowed).To(BeFalse())
					Expect(reason).To(Equal("request is not authenticated"))
				})

				It("should reject unauthorized user", func() {

					req.Request.TLS = &tls.ConnectionState{}
					req.Request.TLS.PeerCertificates = append(req.Request.TLS.PeerCertificates, fakecert)

					result, err := app.generateAccessReview(req)
					Expect(err).ToNot(HaveOccurred())
					result.Status.Allowed = false
					result.Status.Reason = "just because"

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("POST", "/apis/authorization.k8s.io/v1/subjectaccessreviews"),
							ghttp.RespondWithJSONEncoded(http.StatusOK, result),
						),
					)

					allowed, reason, err := app.Authorize(req)
					Expect(err).ToNot(HaveOccurred())
					Expect(allowed).To(BeFalse())
					Expect(reason).To(Equal("just because"))
				})

				It("should allow authorized user", func() {

					req.Request.TLS = &tls.ConnectionState{}
					req.Request.TLS.PeerCertificates = append(req.Request.TLS.PeerCertificates, fakecert)

					result, err := app.generateAccessReview(req)
					Expect(err).ToNot(HaveOccurred())
					result.Status.Allowed = true

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("POST", "/apis/authorization.k8s.io/v1/subjectaccessreviews"),
							ghttp.RespondWithJSONEncoded(http.StatusOK, result),
						),
					)

					allowed, _, err := app.Authorize(req)
					Expect(err).ToNot(HaveOccurred())
					Expect(allowed).To(BeTrue())
				})

				It("should not allow user if auth check fails", func() {

					req.Request.TLS = &tls.ConnectionState{}
					req.Request.TLS.PeerCertificates = append(req.Request.TLS.PeerCertificates, fakecert)

					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("POST", "/apis/authorization.k8s.io/v1/subjectaccessreviews"),
							ghttp.RespondWithJSONEncoded(http.StatusInternalServerError, nil),
						),
					)

					allowed, _, err := app.Authorize(req)
					Expect(err).To(HaveOccurred())
					Expect(allowed).To(BeFalse())
				})
			})

			DescribeTable("should allow all users for info endpoints", func(path string) {
				req.Request.URL.Path = path
				allowed, _, err := app.Authorize(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(allowed).To(BeTrue())
			},
				Entry("root", "/"),
				Entry("apis", "/apis"),
				Entry("group", "/apis/subresources.kubevirt.io"),
				Entry("version", "/apis/subresources.kubevirt.io/version"),
				Entry("healthz", "/apis/subresources.kubevirt.io/healthz"),
				Entry("start profiler", "/apis/subresources.kubevirt.io/start-cluster-profiler"),
				Entry("stop profiler", "/apis/subresources.kubevirt.io/stop-cluster-profiler"),
				Entry("dump profiler", "/apis/subresources.kubevirt.io/dump-cluster-profiler"),
			)

			DescribeTable("should reject all users for unknown endpoint paths", func(path string) {
				req.Request.TLS = &tls.ConnectionState{}
				req.Request.TLS.PeerCertificates = append(req.Request.TLS.PeerCertificates, fakecert)
				req.Request.URL.Path = path
				allowed, _, err := app.Authorize(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(allowed).To(BeFalse())

			},
				Entry("random1", "/apis/subresources.kubevirt.io/v1alpha3/madethisup"),
				Entry("random2", "/1/2/3/4/5/6/7/8/9/0/1/2/3/4/5/6/7/8/9"),
				Entry("no subresource provided", "/apis/subresources.kubevirt.io/v1alpha3/namespaces/default/virtualmachineinstances/testvmi"),
				Entry("invalid resource type", "/apis/subresources.kubevirt.io/v1alpha3/namespaces/default/madeupresource/testvmi/console"),
				Entry("unknown namespaced resource endpoint", "/apis/subresources.kubevirt.io/v1/namespaces/default/madethisup/testvmi/console"),
				Entry("unknown namespaced base resource endpoint", "/apis/subresources.kubevirt.io/v1/namespaces/default/madethisup"),
			)
		})

		AfterEach(func() {
			server.Close()
		})

	})

	DescribeTable("should map verbs", func(httpVerb string, resourceName string, expectedRbacVerb string) {
		Expect(mapHttpVerbToRbacVerb(httpVerb, resourceName)).To(Equal(expectedRbacVerb))
	},
		Entry("http post to create", http.MethodPost, "", "create"),
		Entry("http get with resource to get", http.MethodGet, "foo", "get"),
		Entry("http get without resource to list", http.MethodGet, "", "list"),
		Entry("http put to update", http.MethodPut, "", "update"),
		Entry("http patch to patch", http.MethodPatch, "", "patch"),
		Entry("http delete with reource to delete", http.MethodDelete, "foo", "delete"),
		Entry("http delete without resource to deletecollection", http.MethodDelete, "", "deletecollection"),
	)

})
