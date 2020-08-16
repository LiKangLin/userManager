package main

import (
	"testing"
)

/**
 *由于自定义了init方法，所以需要声明testing的init方法优先执行
 */
var _ = func() bool {
	testing.Init()
	return true
}()

func Test_getUserInfo(t *testing.T) {
	//初始化完成后，功能测试
	existUname := "username138"
	user := getUserInfo(existUname)
	if user.Username != existUname {
		t.Error("用户名不匹配")
	} else {
		t.Log("测试通过")
	}
}

//func Test_editUserInfo(t *testing.T) {
//   // mode 1
//   cnt := editUserInfo("username138", "username138", "", 1)
//   if cnt != 0 && cnt != 1 {
//       t.Error("editUserInfo nickname failed, cnt=", cnt)
//   }
//   // mode 2
//   cnt = editUserInfo("username138", "username138", "www.helloworld.com", 2)
//   if cnt != 0 && cnt != 1 {
//       t.Error("editUserInfo nickname failed, cnt=", cnt)
//   }
//   // mode 3
//   cnt = editUserInfo("username138", "username138", "www.google.com", 3)
//   if cnt != 0 && cnt != 1 {
//       t.Error("editUserInfo nickname failed, cnt=", cnt)
//   }
//}
