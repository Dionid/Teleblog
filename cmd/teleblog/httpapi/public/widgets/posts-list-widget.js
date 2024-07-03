const CROP_TEXT_LENGTH = 500;

window.addEventListener("load", function () {
  const { createApp } = Vue;

  const component = document.querySelector("#posts-list-widget");

  if (!component) {
    return;
  }

  const data = JSON.parse(
    document.getElementById("posts-list-widget-data").textContent
  ).map((item) => {
    return {
      ...item,
      collapsed: item.text.length > CROP_TEXT_LENGTH,
    };
  });

  if (!data) {
    alert("No data found");
    return;
  }

  console.log("data", data);

  createApp({
    data() {
      return {
        loading: false,
        dataById: data.reduce((acc, item) => {
          acc[item.id] = item;
          return acc;
        }, {}),
      };
    },
    methods: {
      cropText(text) {
        console.log("text", text);
        return text.substr(0, CROP_TEXT_LENGTH) + "...";
      },
      expandPostText(postId) {
        this.dataById[postId].collapsed = false;
      },
    },
  }).mount("#posts-list-widget");
});
