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
      collapsed: item.text.length > 200,
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
        return text.substr(0, 200) + "...";
      },
      expandPostText(postId) {
        this.dataById[postId].collapsed = false;
      },
    },
  }).mount("#posts-list-widget");
});
