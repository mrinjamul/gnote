// variables
let username = "";
let notes = [];

noteDocument = document.getElementById("notes");
userDocument = document.getElementById("userinfo");

getData("/user/me").then((data) => {
  username = data.user.user_name;
});

const nonotes = `
      <h1 id="nonewnotes">
        <p>You don't have any notes yet! Go ahead and write some!</p>
      </h1>
      `;
let userinfo = `
      <h1>Welcome, <br />@${username}!</h1>
        <p class="text-primary fs-4">
          You've got <span id="count">${notes.length}</span> note(s) on Gnote.
        </p>
        `;

function submitAndUpdate() {
  createNote = document.getElementById("inputnote");
  if (createNote.innerText == "") {
    return;
  }
  let note = {
    title: "untitled",
    content: createNote.innerText,
  };
  postData("/api/notes", note).then((data) => {
    fetchAndUpdate();
    createNote.innerText = "";
  });
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

async function fetchAndUpdate() {
  const resp = await fetch("/api/notes");
  // if status is not 200, then throw an error
  if (resp.status === 401) {
    logout();
  }
  const data = await resp.json();
  let notes = await data.notes;
  if (notes.length == 0) {
    noteDocument.innerHTML = nonotes;
    userDocument.innerHTML = userinfo;
  } else {
    userDocument.innerHTML = userinfo;

    notesDiv = "";

    for (const note of notes) {
      let date = new Date(note.createdat);
      let noteDiv = `
            <div key=${note.id} class="card isnote">
              <div class="card-body">
                <h5 class="card-title" style="font-size: 0.9em">
                  ${date.toString()}
                </h5>
                <p class="card-text" style="font-size: 1.2em">
                  ${note.content}
                </p>
              </div>
            </div>

            `;

      notesDiv += noteDiv;
    }
    userinfo = `
            <h1>Welcome, <br />@${username}!</h1>
              <p class="text-primary fs-4">
                You've got <span id="count">${notes.length}</span> note(s) on Gnote.
            </p>`;
    userDocument.innerHTML = userinfo;
    noteDocument.innerHTML = notesDiv;
  }
}
fetchAndUpdate();

// Add event listener to keypress
document.addEventListener("keypress", function (event) {
  if ((event.keyCode == 10 || event.keyCode == 13) && event.ctrlKey) {
    submitAndUpdate();
  }
});

function logout() {
  postData("/auth/logout", {}).then((data) => {
    window.location.href = "/";
  });
}

async function getData(url = "") {
  // Default options are marked with *
  const response = await fetch(url, {
    method: "GET", // *GET, POST, PUT, DELETE, etc.
    mode: "cors", // no-cors, *cors, same-origin
    cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
    credentials: "same-origin", // include, *same-origin, omit
    headers: {
      "Content-Type": "application/json",
      // 'Content-Type': 'application/x-www-form-urlencoded',
    },
    redirect: "follow", // manual, *follow, error
    referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    // body data type must match "Content-Type" header
  });
  return response.json(); // parses JSON response into native JavaScript objects
}
