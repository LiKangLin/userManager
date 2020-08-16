package main

import (
	"testing"
)

func Test_getTableName(t *testing.T) {
	tab := getTableName("username138")
	if tab != "userinfo_tab_0" {
		t.Error("table name error:", tab)
	}
}

func Test_getDbUserInfo(t *testing.T) {
	username := "username138"
	user, err := getDbUserInfo(username)
	if err != nil {
		t.Error("getDbUserInfo failed, ", err.Error())
	} else if username != user.Username {
		t.Error("getDbUserInfo not match, name: ", user.Username)
	}
}

func Test_updateDbNickname(t *testing.T) {
	cnt := updateDbNickname("username138", "godlike")
	if cnt != 1 && cnt != 0 {
		t.Error("updateDbNickname failed, affected_rows is not 1 or 0, cnt=", cnt)
	}
}

func Test_updateDbHeadurl(t *testing.T) {
	cnt := updateDbHeadurl("username138", "www.google.cn")
	if cnt != 1 && cnt != 0 {
		t.Error("updateDbHeadurl failed, affected_rows is not 1 or 0, cnt=", cnt)
	}
}

func Test_updateDbUserinfo(t *testing.T) {
	cnt := updateDbUserinfo("username138", "godlike", "www.google.cn")
	if cnt != 1 && cnt != 0 {
		t.Error("updateDbUserinfo failed, affected_rows is not 1 or 0, cnt=", cnt)
	}
}
