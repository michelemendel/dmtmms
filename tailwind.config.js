/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './public/**/*.{html,js}',
    './view/**/*_templ.go',
    './view/**/*_templ.txt',
    // "./node_modules/tw-elements/dist/js/**/*.js"
  ],
  theme: {
    // themes: false, // false: only light + dark | true: all themes | array: specific themes like this ["light", "dark", "cupcake"]
    // darkTheme: "light", // name of one of the included themes for dark mode
    // base: true, // applies background color and foreground color for root element by default
    // styled: true, // include daisyUI colors and design decisions for all components
    // utils: true, // adds responsive and modifier utility classes
    // prefix: "", // prefix for daisyUI classnames (components, modifiers and responsive class names. Not colors)
    // logs: true, // Shows info about daisyUI version and used config in the console when building your CSS
    // themeRoot: ":root",
  },
  plugins: [
    require('@tailwindcss/forms'),
    // require("daisyui"),
    // require("tw-elements/dist/plugin.cjs"),
  ],
  darkMode: "false",
  // daisyui: {
  //   themes: ["light"],
  // },
}