/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './public/**/*.{html,js}', 
    './view/**/*_templ.go',
    './view/**/*_templ.txt',
  ],
  theme: {
  },
  plugins: [
    // require('@tailwindcss/forms'),    
  ],
}