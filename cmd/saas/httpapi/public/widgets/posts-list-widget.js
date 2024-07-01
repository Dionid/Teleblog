window.addEventListener("load", function () {
  const { createApp } = Vue;

  const component = document.querySelector("#posts-list-widget");

  if (!component) {
    return;
  }

  const data = JSON.parse(
    document.getElementById("presentation-data").textContent
  );

  if (!data) {
    alert("No data found");
    return;
  }

  createApp({
    data() {
      return {
        loading: false,
      };
    },
  }).mount("#posts-list-widget");
});
