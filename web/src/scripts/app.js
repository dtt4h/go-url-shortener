const formElement = document.getElementById("mainForm");
const urlString = document.getElementById("urlString");
const resultContainer = document.getElementById("resultContainer");

formElement.addEventListener("submit", (evt) => {
  evt.preventDefault()

  const urlResult = document.querySelector(".urlResult");

  urlResult.textContent = `${urlString.value}`
})