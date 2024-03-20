/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './public/**/*.{html,js}',
    './view/**/*_templ.go',
    './view/**/*_templ.txt',
    './view/**/formelements.go',
  ],
  theme: {
    extend: {
      colors: {
        'i-blue': '#0038b8',
      },
    }
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