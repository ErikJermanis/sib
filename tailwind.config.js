/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.templ"],
  theme: {
    extend: {
      colors: {
        night: "#0D0A0B",
        charcoal: "#454955",
        magnolia: "#F3EFF5",
        magnoliaDark: "#E8E4EB",
        magnoliaAccent: "#DDCFE7",
        appleGreen: "#72B01D",
        officeGreen: "#3F7D20",
        danger: "#A90916",
      },
    },
  },
  ignoreFiles: ["node_modules/**", "bin/**"],
  plugins: [],
};
