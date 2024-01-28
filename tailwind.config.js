/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './public/**/*.{html,js}', 
    './view/**/*_templ.go',
  ],
  theme: {
  },
  plugins: [
    // require('@tailwindcss/forms'),    
  ],
}