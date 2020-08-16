function login() {
    var username = document.getElementById("username")
    var password = document.getElementById("password")

    if (username.value == "") {
        alert("请输入用户名")
    } else if (password.value == "") {
        alert("请输入密码")
    }
    if (!username.value.match(/^\S{2,20}$/)) {
        console.log("get focus")
        username.className = 'userRed';
        username.focus();
        return;
    }

    var xhr = new XMLHttpRequest();
    xhr.open('post', 'http://127.0.0.1:8080/login')
    xhr.setRequestHeader("Content-type","application/x-www-form-urlencoded")
    //xhr.withCredentials=true
    xhr.send('username=' + username.value + "&passwd=" + password.value)
    console.log('username=' + username.value + "&passwd=" + password.value)
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            console.log(xhr.responseText)
            var json = eval("("+xhr.responseText+")");
            console.log(json.code)
            console.log(json.data)
            //调试使用
            //alert(xhr.responseText)
            if (json.code == 0) {
                //setCookie('token', json.data.token)
                //console.log(json.data.token)
                //cookie.attr("token", json.data.token)
                window.location.href = "http://127.0.0.1:8080/static/index.html?name=" + username.value+"&token="+json.data.token
                window.event.returnValue = false
            } else {
                alert("账号或密码错误。")
            }
        }
    }
}

// 封装cookie

/**
 * 设置cookie
 * @param {*} name cookie名称
 * @param {*} value cookie值
 * @param {*} iDay 过期时间（天数）
 */
function setCookie(name, value, iDay)
{
    var oDate=new Date();
    oDate.setDate(oDate.getDate()+iDay);

    document.cookie=name+'='+value+';expires='+oDate;
}

function getCookie(name)
{
    var arr=document.cookie.split('; ');

    for(var i=0;i<arr.length;i++)
    {
        var arr2=arr[i].split('=');

        if(arr2[0]==name)
        {
            return arr2[1];
        }
    }

    return '';
}

function removeCookie(name)
{
    setCookie(name, 1, -1);
}
