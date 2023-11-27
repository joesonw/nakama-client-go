package tests

import (
	"context"

	"github.com/heroiclabs/nakama-common/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Friend Tests", func() {
	It("should add friend, then list", func() {
		id1 := generateID()
		id2 := generateID()
		client := newNakamaHTTPClient()

		sess1, err := client.AuthenticateCustom(context.Background(), id1)
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), id2)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(sess1.AddFriends(context.Background(), &api.AddFriendsRequest{
			Ids: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		res, err := sess1.ListFriends(context.Background(), &api.ListFriendsRequest{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.Friends).To(HaveLen(1))
		Expect(res.Friends[0].State.GetValue()).To(BeEquivalentTo(api.Friend_INVITE_SENT))
	})

	It("should receive friend invite, then list", func() {
		id1 := generateID()
		id2 := generateID()
		client := newNakamaHTTPClient()

		sess1, err := client.AuthenticateCustom(context.Background(), id1)
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), id2)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(sess1.AddFriends(context.Background(), &api.AddFriendsRequest{
			Ids: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		res, err := sess2.ListFriends(context.Background(), &api.ListFriendsRequest{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.Friends).To(HaveLen(1))
		Expect(res.Friends[0].State.GetValue()).To(BeEquivalentTo(api.Friend_INVITE_RECEIVED))
	})

	It("should block friend, then list", func() {
		id1 := generateID()
		id2 := generateID()
		client := newNakamaHTTPClient()

		sess1, err := client.AuthenticateCustom(context.Background(), id1)
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), id2)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(sess1.BlockFriends(context.Background(), &api.BlockFriendsRequest{
			Ids: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		res, err := sess1.ListFriends(context.Background(), &api.ListFriendsRequest{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.Friends).To(HaveLen(1))
		Expect(res.Friends[0].State.GetValue()).To(BeEquivalentTo(api.Friend_BLOCKED))
	})

	It("should add friend, accept, then list", func() {
		id1 := generateID()
		id2 := generateID()
		client := newNakamaHTTPClient()

		sess1, err := client.AuthenticateCustom(context.Background(), id1)
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), id2)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(sess1.AddFriends(context.Background(), &api.AddFriendsRequest{
			Ids: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess2.AddFriends(context.Background(), &api.AddFriendsRequest{
			Ids: []string{sess1.UserID()},
		})).ShouldNot(HaveOccurred())

		res, err := sess1.ListFriends(context.Background(), &api.ListFriendsRequest{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.Friends).To(HaveLen(1))
		Expect(res.Friends[0].State.GetValue()).To(BeEquivalentTo(api.Friend_FRIEND))
	})

	It("should add friend, reject, then list", func() {
		id1 := generateID()
		id2 := generateID()
		client := newNakamaHTTPClient()

		sess1, err := client.AuthenticateCustom(context.Background(), id1)
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), id2)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(sess1.AddFriends(context.Background(), &api.AddFriendsRequest{
			Ids: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess2.DeleteFriends(context.Background(), &api.DeleteFriendsRequest{
			Ids: []string{sess1.UserID()},
		})).ShouldNot(HaveOccurred())

		res, err := sess1.ListFriends(context.Background(), &api.ListFriendsRequest{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.Friends).To(HaveLen(0))
	})

	// It("should add friend, accept, then list", func() {
	//	id1 := generateID()
	//	id2 := generateID()
	//	client := newNakamaHTTPClient()
	//
	//	sess1, err := client.AuthenticateCustom(context.Background(), id1)
	//	Expect(err).ShouldNot(HaveOccurred())
	//	sess2, err := client.AuthenticateFacebookInstantGame(context.Background(), createFacebookInstantGameAuthToken(id2))
	//	Expect(err).ShouldNot(HaveOccurred())
	//
	//	Expect(sess1.AddFriends(context.Background(), &api.AddFriendsRequest{
	//		Ids: []string{sess2.UserID()},
	//	})).ShouldNot(HaveOccurred())
	//	res, err := sess1.ListFriends(context.Background(), &api.ListFriendsRequest{})
	//	Expect(err).ShouldNot(HaveOccurred())
	//	Expect(res).NotTo(BeNil())
	//	Expect(res.Friends).To(HaveLen(1))
	//	Expect(res.Friends[0].User.GetFacebookInstantGameId()).To(Equal(id2))
	// })
})
