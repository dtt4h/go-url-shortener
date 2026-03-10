const formElement = document.getElementById("mainForm");
const urlString = document.getElementById("urlString");
const resultContainer = document.getElementById("resultContainer");

const urlResult = document.querySelector(".urlResult");

formElement.addEventListener("submit", (evt) => {
  evt.preventDefault();
  shortResultContainer();
  const newPost = {url: urlString.value}; 
  createShortUrl(newPost);
})

const shortResultContainer = () => {
  resultContainer.classList.add("is-active");
}

const createShortUrl = (newPost) => {
  return fetch("http://127.0.0.1:8080/api/v1/shorten", {
      method: "POST",
      body: JSON.stringify(newPost),
      headers: {
        "Content-type": "application/json",
      },
    })
    .then((res) => res.json())
    .then((data)=> urlResult.textContent = "click.ru/"+(data.short_url.split("/")).pop())//заменить клилкюрл на нашу тему
}
