// refresh.js

let duration = 275;
refreshToken(duration);

function refreshToken(duration) {
  // if (checkCookie("token")) {
  //   console.log("cookie exists");
  // }

  // Run without check
  GETRefresh();
  setInterval(GETRefresh(), duration * 1000);
}

function GETRefresh() {
  postData("/auth/refresh", {}).then((data) => {
    // console.log("token refreshed");
  });
}
// check if cookie exists or not
function checkCookie(name = "") {
  var value = getCookie(name);
  if (value == null) {
    return false;
  } else {
    return true;
  }
}

// getCookie returns cookie value
function getCookie(name) {
  var dc = document.cookie;
  var prefix = name + "=";
  var begin = dc.indexOf("; " + prefix);
  if (begin == -1) {
    begin = dc.indexOf(prefix);
    if (begin != 0) return null;
  } else {
    begin += 2;
    var end = document.cookie.indexOf(";", begin);
    if (end == -1) {
      end = dc.length;
    }
  }
  // because unescape has been deprecated, replaced with decodeURI
  //return unescape(dc.substring(begin + prefix.length, end));
  return decodeURI(dc.substring(begin + prefix.length, end));
}
