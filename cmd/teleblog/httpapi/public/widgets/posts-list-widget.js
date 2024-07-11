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
      const query = new URLSearchParams(window.location.search);

      return {
        loading: false,
        dataById: data.reduce((acc, item) => {
          acc[item.id] = item;
          return acc;
        }, {}),
        searchString: query.get("search") || "",
        tag: query.get("tag") || "_",
      };
    },
    watch: {
      tag() {
        this.search();
      },
    },
    methods: {
      cropText(text) {
        return text.substr(0, CROP_TEXT_LENGTH) + "...";
      },
      expandPostText(postId) {
        this.dataById[postId].collapsed = false;
      },
      search() {
        const query = new URLSearchParams(window.location.search);

        query.set("page", 1);
        query.set("search", this.searchString);
        if (this.tag !== "_") {
          query.set("tag", this.tag);
        }

        window.location = `?${query.toString()}`;
      },
      setPage(pageNum, event) {
        if (event) {
          event.preventDefault();
        }

        const query = new URLSearchParams(window.location.search);

        query.set("page", pageNum);

        window.location = `?${query.toString()}`;
      },
    },
  }).mount("#posts-list-widget");
});
