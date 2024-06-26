window.addEventListener("load", function () {
  const { createApp } = Vue;
  createApp({
    data() {
      return {
        email: "",
        error: "",
        loading: false,
        success: "",
      };
    },
    watch: {
      email() {
        this.error = "";
      },
    },
    methods: {
      async request() {
        if (this.email === "") {
          this.error = "Email is required";
          return;
        }

        this.loading = true;

        const data = {
          email: this.email,
        };

        const response = await fetch(
          "/api/collections/users/request-password-reset",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
          }
        );

        this.loading = false;

        if (!response.ok) {
          const respJson = await response.json();
          this.error = respJson.message;
          return;
        }

        this.success = "Password reset link sent to your email";
      },
    },
  }).mount("#reset-password-form-component");
});
