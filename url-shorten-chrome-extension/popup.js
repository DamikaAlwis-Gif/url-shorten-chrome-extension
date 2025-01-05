document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("shorten-form");
  const urlInput = document.getElementById("url-input");
  const customInput = document.getElementById("custom-input");
  const resultDiv = document.getElementById("result");
  const shortenedUrl = document.getElementById("shortened-url");
  const copyButton = document.getElementById("copy-button");
  const errorMessage = document.getElementById("error-message");

  form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const url = urlInput.value;
    const customUrl = customInput.value;

    // Clear previous error messages
    errorMessage.textContent = "";
    errorMessage.classList.add("hidden");

    // URL Validation
    if (!isValidUrl(url)) {
      // alert("Please enter a valid URL.");
      errorMessage.textContent = "Please enter a valid URL.";
      errorMessage.classList.remove("hidden");
      return;
    }

    try {
      // Replace with your backend API endpoint
      const response = await fetch("http://localhost:5000/shorten", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url, "short" : customUrl }),
      });

      if (!response.ok) {
        if (response.status == 400){
          data = await response.json();
          error = data.error
          if (error == "short key already exists"){
            errorMessage.textContent = "Short key already exists."
            errorMessage.classList.remove("hidden");
          }
        }
        
        
        console.error("Failed to shorten the URL, status:", response.status);
        throw new Error("Failed to shorten the URL");
      }

      const data = await response.json();
      console.log("Shortened URL received:", data.short);
      shortenedUrl.href = data.short;
      shortenedUrl.textContent = data.short;
      resultDiv.classList.remove("hidden");
    } catch (error) {
      console.error("Error shortening the URL:", error.message);
      errorMessage.textContent = "Error: " + error.message;
      errorMessage.classList.remove("hidden");
    }
  });

  // Copy to Clipboard functionality
  copyButton.addEventListener("click", () => {
    const range = document.createRange();
    range.selectNode(shortenedUrl);
    window.getSelection().removeAllRanges();
    window.getSelection().addRange(range);
    document.execCommand("copy");
    alert("URL copied to clipboard!");
  });

  // Function to validate the URL
  function isValidUrl(url) {
    try {
      const parsedUrl = new URL(url);

      //  allow only the standard web protocols
      const allowedProtocols = ["http:", "https:"]
      return allowedProtocols.includes(parsedUrl.protocol)
      
    } catch (e) {
      return false;
    }
  }
});
