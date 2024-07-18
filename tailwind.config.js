/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/components/**/*.templ",
    "./node_modules/flowbite/**/*.js",
  ],
  theme: {
    extend: {},
  },
  plugins: [require("flowbite/plugin")],
};
