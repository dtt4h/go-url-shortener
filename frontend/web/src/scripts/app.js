const formElement = document.getElementById("mainForm");
const urlString = document.getElementById("urlString");
const resultContainer = document.getElementById("resultContainer");
const copyButton = document.getElementById("copyButton");
const qrCode = document.getElementById("qrCode");
const link = document.getElementById("stringLink");
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



const API_BASE = window.VERCEL_URL || window.location.origin;

const createShortUrl = (newPost) => { 
  return fetch(`${API_BASE}/api/v1/shorten`, {
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
      currentShortUrl = data.short_url;
      urlResult.textContent = currentShortUrl;
      qrCode.src= (data.short_url + "/qr").trim();
      stringLink.href=currentShortUrl;
      shortResultContainer();
    }) //заменить клилкюрл на нашу тему
    .catch((err) => {
      console.log("Create short url failed:", err);
      currentShortUrl = ""
      urlResult.textContent = ""
    })
}
