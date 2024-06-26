window.addEventListener("load", function () {
  const { createApp } = Vue;
  createApp({
    data() {
      return {
        email: "",
        password: "",
        error: "",
        loading: false,
      };
    },
    watch: {
      password() {
        this.error = "";
      },
      email() {
        this.error = "";
      },
    },
    methods: {
      async signUp() {
        if (this.email === "" || this.password === "") {
          this.error = "Email and password are required";
          return;
        }

        this.loading = true;

        const data = {
          email: this.email,
          emailVisibility: true,
          password: this.password,
          passwordConfirm: this.password,
        };

        const response = await fetch("/api/collections/users/records", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(data),
        });

        this.loading = false;

        if (!response.ok) {
          const respJson = await response.json();
          this.error = respJson.message;
          return;
        }

        window.location.href = "/auth/sign-in";
      },
    },
  }).mount("#sign-up-form-component");
});
