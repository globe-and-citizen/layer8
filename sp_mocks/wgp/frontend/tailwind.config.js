/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {},
  },
  plugins: [require('@tailwindcss/typography'), require("daisyui")],

  daisyui: {
    themes: [
      {
        custom: {
         // ...require("daisyui/src/theming/themes")["light"],
          primary: "#4caf50",
          secondary: "#86EE89",
          accent: "#43FB49",
          // "neutral": "#3d4451",
          "base-100": "#ffffff",
        },
      },
      "light",
      "dark",
      "cupcake",
    ],
  },
};
