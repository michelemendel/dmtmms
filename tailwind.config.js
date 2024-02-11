/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './public/**/*.{html,js}',
    './view/**/*_templ.go',
    './view/**/*_templ.txt',
    "./node_modules/tw-elements/dist/js/**/*.js"
  ],
  theme: {
  },
  plugins: [
    require('@tailwindcss/forms'),
    require("tw-elements/dist/plugin.cjs"),
  ],
  darkMode: "class"
}