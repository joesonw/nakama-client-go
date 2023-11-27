package tests

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nakama_client_go "github.com/joesonw/nakama-client-go"
)

var _ = Describe("Authenticate Tests", func() {
	It("should authenticate with email", func() {
		sess, err := newNakamaHTTPClient().AuthenticateEmail(context.Background(), generateID()+"@example.com", generateID())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess).NotTo(BeNil())
		Expect(sess.Token()).NotTo(BeEmpty())
	})

	It("should authenticate with device id", func() {
		sess, err := newNakamaHTTPClient().AuthenticateDevice(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess).NotTo(BeNil())
		Expect(sess.Token()).NotTo(BeEmpty())
	})

	It("should authenticate with custom id", func() {
		sess, err := newNakamaHTTPClient().AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess).NotTo(BeNil())
		Expect(sess.Token()).NotTo(BeEmpty())
	})

	It("should authenticate with user variables", func() {
		sess, err := newNakamaHTTPClient().AuthenticateCustom(context.Background(), generateID(), nakama_client_go.AuthenticateOption{
			Vars: map[string]string{
				"testString": "testValue",
			},
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess).NotTo(BeNil())
		Expect(sess.Token()).NotTo(BeEmpty())
		Expect(sess.Var("testString")).To(Equal("testValue"))
	})

	It("should fail to authenticate with new custom id", func() {
		_, err := newNakamaHTTPClient().AuthenticateCustom(context.Background(), generateID(), nakama_client_go.AuthenticateOption{
			Create: nakama_client_go.False(),
		})
		Expect(err).Should(HaveOccurred())
	})

	It("should authenticate with custom id twice", func() {
		id := generateID()
		_, err := newNakamaHTTPClient().AuthenticateCustom(context.Background(), id)
		Expect(err).ShouldNot(HaveOccurred())
		sess, err := newNakamaHTTPClient().AuthenticateCustom(context.Background(), id)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess).NotTo(BeNil())
		Expect(sess.Token()).NotTo(BeEmpty())
	})

	It("should fail to authenticate with custom id", func() {
		_, err := newNakamaHTTPClient().AuthenticateCustom(context.Background(), "")
		Expect(err).Should(HaveOccurred())
	})

	// It("should authenticate with facebook instant game", func() {
	//	sess, err := newNakamaHTTPClient().AuthenticateFacebookInstantGame(context.Background(), createFacebookInstantGameAuthToken("a_player_id"))
	//	Expect(err).ShouldNot(HaveOccurred())
	//	Expect(sess).NotTo(BeNil())
	//	Expect(sess.Token()).NotTo(BeEmpty())
	// })
})
