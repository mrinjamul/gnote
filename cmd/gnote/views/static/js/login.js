const pwd = document.getElementById("pass");
const vpwd = document.getElementById("vpass");
const form = document.getElementById("signup-form");

function verifyPassword() {
  var pass = pwd.value;
  var vpass = vpwd.value;
  if (pass == vpass) {
    form.submit();
  } else {
    alert("The passwords don't match");
  }
}

async function postData(url = "", data = {}) {
  // Default options are marked with *
  const response = await fetch(url, {
    method: "POST", // *GET, POST, PUT, DELETE, etc.
    mode: "cors", // no-cors, *cors, same-origin
    cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
    credentials: "same-origin", // include, *same-origin, omit
    headers: {
      "Content-Type": "application/json",
      // 'Content-Type': 'application/x-www-form-urlencoded',
    },
    redirect: "follow", // manual, *follow, error
    referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: JSON.stringify(data), // body data type must match "Content-Type" header
  });
  return response.json(); // parses JSON response into native JavaScript objects
}

// checkEmpty check if username or password is empty
function checkEmpty(username, password) {
  if (username == "" || password == "") {
    return false;
  } else {
    return true;
  }
}
