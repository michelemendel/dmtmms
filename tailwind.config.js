/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './public/**/*.{html,js}',
    './view/**/*_templ.go',
    './view/**/*_templ.txt',
  ],
  theme: {
  },
  variants: {
    extend: {
      borderStyle: ['responsive', 'hover'],
      borderWidth: ['responsive', 'hover'],
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
  darkMode: "false",
}