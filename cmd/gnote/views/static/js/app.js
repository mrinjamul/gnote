// variables
let username = "";
let notes = [];

noteDocument = document.getElementById("notes");
userDocument = document.getElementById("userinfo");

const nonotes = `
      <h1 id="nonewnotes">
        <p>You don't have any notes yet! Go ahead and write some!</p>
      </h1>
      `;
let userinfo;

getData("/user/me").then((data) => {
  username = data.user.username;
  userinfo = `
      <h1>Welcome, <br />@${username}!</h1>
        <p class="text-primary fs-4">
          You've got <span id="count">${notes.length}</span> note(s) on Gnote.
        </p>
      `;
});

function submitAndUpdate() {
  let noteTitle = document.getElementById("inputtitle");
  let createNote = document.getElementById("inputnote");
  if (createNote.innerText == "") {
    return;
  }
  title = noteTitle.innerHTML != "" ? noteTitle.innerHTML : "untitled";
  let note = {
    title: title,
    content: createNote.innerText,
  };
  postData("/api/notes", note).then((data) => {
    fetchAndUpdate();
    createNote.innerText = "";
  });
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
      let date = new Date(note.created_at);
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

      notesDiv = noteDiv + notesDiv;
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
