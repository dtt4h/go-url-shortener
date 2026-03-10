const formElement = document.getElementById("mainForm");
const urlString = document.getElementById("urlString");
const resultContainer = document.getElementById("resultContainer");
const copyButton = document.getElementById("copyButton");
let currentShortUrl;

const urlResult = document.querySelector(".urlResult");

formElement.addEventListener("submit", (evt) => {
  evt.preventDefault();
  const newPost = {url: urlString.value}; 
  createShortUrl(newPost);
})

copyButton.addEventListener("click", copyUrl);

const shortResultContainer = () => {
  resultContainer.classList.add("is-active");
}

async function copyUrl() {
  try{
    await navigator.clipboard.writeText(currentShortUrl);
    console.log("Done!");
  } catch (err) {
    console.log("Error by copy", err);
  }
}


const createShortUrl = (newPost) => { 
  return fetch("http://127.0.0.1:8080/api/v1/shorten", {
      method: "POST",
      body: JSON.stringify(newPost),
      headers: {
        "Content-type": "application/json",
      },
    })
    .then((res) => {
      if (res.ok){
        return res.json()
      }
      else{
        throw new Error("Failed to get link from server")
      }
    })
    .then((data)=> {
      currentShortUrl = "click.ru/"+(data.short_url.split("/")).pop();
      urlResult.textContent = currentShortUrl;
      shortResultContainer();
    }) //заменить клилкюрл на нашу тему
    .catch((err) => {
      console.log("Create short url failed:", err);
      currentShortUrl = ""
      urlResult.textContent = ""
    })
}
