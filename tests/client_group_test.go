package tests

import (
	"context"

	"github.com/heroiclabs/nakama-common/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Group Tests", func() {
	It("should create group", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err)
		group, err := sess.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(group).NotTo(BeNil())
		Expect(group.Name).To(Equal(groupName))
		Expect(group.EdgeCount).To(BeEquivalentTo(1))
		Expect(group.MaxCount).To(BeEquivalentTo(100))
		Expect(group.Metadata).To(Equal("{}"))
		Expect(group.GetOpen().GetValue()).To(BeTrue())
	})

	It("should create then list group", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		_, err = sess.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		list, err := sess.ListUserGroups(context.Background(), &api.ListUserGroupsRequest{
			UserId: sess.UserID(),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetUserGroups()).To(HaveLen(1))
		Expect(list.GetUserGroups()[0].GetGroup().Name).To(Equal(groupName))
		Expect(list.GetUserGroups()[0].GetGroup().EdgeCount).To(BeEquivalentTo(1))
		Expect(list.GetUserGroups()[0].GetGroup().MaxCount).To(BeEquivalentTo(100))
		Expect(list.GetUserGroups()[0].GetGroup().Metadata).To(Equal("{}"))
		Expect(list.GetUserGroups()[0].GetGroup().GetOpen().GetValue()).To(BeTrue())
		Expect(list.GetUserGroups()[0].GetState().GetValue()).To(BeEquivalentTo(api.UserGroupList_UserGroup_SUPERADMIN))
	})

	It("should create, update then list group", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		groupName2 := generateID()
		sess, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess.UpdateGroup(context.Background(), &api.UpdateGroupRequest{
			GroupId: group.Id,
			Name:    wrapperspb.String(groupName2),
		})).ShouldNot(HaveOccurred())
		list, err := sess.ListUserGroups(context.Background(), &api.ListUserGroupsRequest{
			UserId: sess.UserID(),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetUserGroups()).To(HaveLen(1))
		Expect(list.GetUserGroups()[0].GetGroup().Name).To(Equal(groupName2))
		Expect(list.GetUserGroups()[0].GetGroup().EdgeCount).To(BeEquivalentTo(1))
		Expect(list.GetUserGroups()[0].GetGroup().MaxCount).To(BeEquivalentTo(100))
		Expect(list.GetUserGroups()[0].GetGroup().Metadata).To(Equal("{}"))
		Expect(list.GetUserGroups()[0].GetGroup().GetOpen().GetValue()).To(BeTrue())
		Expect(list.GetUserGroups()[0].GetState().GetValue()).To(BeEquivalentTo(api.UserGroupList_UserGroup_SUPERADMIN))
	})

	It("should create, delete then list group", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess.DeleteGroup(context.Background(), &api.DeleteGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		list, err := sess.ListUserGroups(context.Background(), &api.ListUserGroupsRequest{
			UserId: sess.UserID(),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetUserGroups()).To(HaveLen(0))
	})

	It("should create then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		list, err := sess.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(1))
		Expect(list.GetGroupUsers()[0].GetState().GetValue()).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
	})

	It("should create, add, then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(2))
		user1 := list.GetGroupUsers()[0].GetState().GetValue()
		user2 := list.GetGroupUsers()[1].GetState().GetValue()
		if user1 < user2 {
			Expect(user1).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
			Expect(user2).To(BeEquivalentTo(api.GroupUserList_GroupUser_MEMBER))
		} else {
			Expect(user1).To(BeEquivalentTo(api.GroupUserList_GroupUser_MEMBER))
			Expect(user2).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
		}
	})

	It("should create, add, kick, then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess1.KickGroupUsers(context.Background(), &api.KickGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		}))
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(1))
		Expect(list.GetGroupUsers()[0].GetState().GetValue()).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
	})

	It("should create closed, request to join, then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: false,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess2.JoinGroup(context.Background(), &api.JoinGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(2))
		user1 := list.GetGroupUsers()[0].GetState().GetValue()
		user2 := list.GetGroupUsers()[1].GetState().GetValue()
		if user1 < user2 {
			Expect(user1).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
			Expect(user2).To(BeEquivalentTo(api.GroupUserList_GroupUser_JOIN_REQUEST))
		} else {
			Expect(user1).To(BeEquivalentTo(api.GroupUserList_GroupUser_JOIN_REQUEST))
			Expect(user2).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
		}
	})

	It("should create open, join, then list group users and count affected", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess2.JoinGroup(context.Background(), &api.JoinGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListUserGroups(context.Background(), &api.ListUserGroupsRequest{
			UserId: sess1.UserID(),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetUserGroups()).To(HaveLen(1))
		Expect(list.GetUserGroups()[0].GetGroup().GetEdgeCount()).To(BeEquivalentTo(2))
	})

	It("should create closed, request to join, then list user groups and see count unaffected", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: false,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess2.JoinGroup(context.Background(), &api.JoinGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListUserGroups(context.Background(), &api.ListUserGroupsRequest{
			UserId: sess1.UserID(),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetUserGroups()).To(HaveLen(1))
		Expect(list.GetUserGroups()[0].GetGroup().GetEdgeCount()).To(BeEquivalentTo(1))
	})

	It("should create closed, request to join, add, then list user groups and see count affected", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: false,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess2.JoinGroup(context.Background(), &api.JoinGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListUserGroups(context.Background(), &api.ListUserGroupsRequest{
			UserId: sess1.UserID(),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetUserGroups()).To(HaveLen(1))
		Expect(list.GetUserGroups()[0].GetGroup().GetEdgeCount()).To(BeEquivalentTo(2))
	})

	It("should create open, join, leave, then list user groups and see count affected", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess2.JoinGroup(context.Background(), &api.JoinGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		Expect(sess2.LeaveGroup(context.Background(), &api.LeaveGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListUserGroups(context.Background(), &api.ListUserGroupsRequest{
			UserId: sess1.UserID(),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetUserGroups()).To(HaveLen(1))
		Expect(list.GetUserGroups()[0].GetGroup().GetEdgeCount()).To(BeEquivalentTo(1))
	})

	It("should create closed, request to join, leave, then list user groups and see count unaffected", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: false,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess2.JoinGroup(context.Background(), &api.JoinGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		Expect(sess2.LeaveGroup(context.Background(), &api.LeaveGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListUserGroups(context.Background(), &api.ListUserGroupsRequest{
			UserId: sess1.UserID(),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetUserGroups()).To(HaveLen(1))
		Expect(list.GetUserGroups()[0].GetGroup().GetEdgeCount()).To(BeEquivalentTo(1))
	})

	It("should create closed, request to join, add, then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess2.JoinGroup(context.Background(), &api.JoinGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(2))
		user1 := list.GetGroupUsers()[0].GetState().GetValue()
		user2 := list.GetGroupUsers()[1].GetState().GetValue()
		if user1 < user2 {
			Expect(user1).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
			Expect(user2).To(BeEquivalentTo(api.GroupUserList_GroupUser_MEMBER))
		} else {
			Expect(user1).To(BeEquivalentTo(api.GroupUserList_GroupUser_MEMBER))
			Expect(user2).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
		}
	})

	It("should create closed, request to join, kick, then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess2.JoinGroup(context.Background(), &api.JoinGroupRequest{
			GroupId: group.Id,
		})).ShouldNot(HaveOccurred())
		Expect(sess1.KickGroupUsers(context.Background(), &api.KickGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(1))
		Expect(list.GetGroupUsers()[0].GetState().GetValue()).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
	})

	It("should create, add, promote, then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess1.PromoteGroupUsers(context.Background(), &api.PromoteGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(2))
		user1 := list.GetGroupUsers()[0].GetState().GetValue()
		user2 := list.GetGroupUsers()[1].GetState().GetValue()
		if user1 < user2 {
			Expect(user1).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
			Expect(user2).To(BeEquivalentTo(api.GroupUserList_GroupUser_ADMIN))
		} else {
			Expect(user1).To(BeEquivalentTo(api.GroupUserList_GroupUser_ADMIN))
			Expect(user2).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
		}
	})

	It("should create, add, promote twice, then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess1.PromoteGroupUsers(context.Background(), &api.PromoteGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess1.PromoteGroupUsers(context.Background(), &api.PromoteGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(2))
		Expect(list.GetGroupUsers()[0].GetState().GetValue()).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
		Expect(list.GetGroupUsers()[1].GetState().GetValue()).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
	})

	It("should create, add, promote, add, promote then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess3, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess1.PromoteGroupUsers(context.Background(), &api.PromoteGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess2.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess3.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess2.PromoteGroupUsers(context.Background(), &api.PromoteGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess3.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(3))
	})

	It("should create add, ban, then list group users", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess1.BanGroupUsers(context.Background(), &api.BanGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(1))
	})

	It("should create group and fail to demote superadmin", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess.DemoteGroupUsers(context.Background(), &api.DemoteGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(1))
		Expect(list.GetGroupUsers()[0].GetState().GetValue()).To(BeEquivalentTo(api.GroupUserList_GroupUser_SUPERADMIN))
	})

	It("should create group, add user and promote and demote them", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess1.PromoteGroupUsers(context.Background(), &api.PromoteGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess1.DemoteGroupUsers(context.Background(), &api.DemoteGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
			State:   wrapperspb.Int32(int32(api.GroupUserList_GroupUser_MEMBER)),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(1))
	})

	It("should create group, add user and fail to demote lowest-ranking user", func() {
		client := newNakamaHTTPClient()
		groupName := generateID()
		sess1, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		sess2, err := client.AuthenticateCustom(context.Background(), generateID())
		Expect(err).ShouldNot(HaveOccurred())
		group, err := sess1.CreateGroup(context.Background(), &api.CreateGroupRequest{
			Name: groupName,
			Open: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sess1.AddGroupUsers(context.Background(), &api.AddGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		Expect(sess1.DemoteGroupUsers(context.Background(), &api.DemoteGroupUsersRequest{
			GroupId: group.Id,
			UserIds: []string{sess2.UserID()},
		})).ShouldNot(HaveOccurred())
		list, err := sess1.ListGroupUsers(context.Background(), &api.ListGroupUsersRequest{
			GroupId: group.Id,
			State:   wrapperspb.Int32(int32(api.GroupUserList_GroupUser_MEMBER)),
		})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(list.GetGroupUsers()).To(HaveLen(1))
	})
})
