/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
	'./views/**/*.templ',
	'./views/**/*.go',
  ],
  theme: {
    screens: {
      sm: '480px',
      md: '768px',
      lg: '976px',
      xl: '1440px',
    },
    colors: {
	    'mywhite': '#ebf2fa',
	    'dark': '#25283d',
	    'green': '#297373',
	    'yellow': '#bdd358',
	    'violet': '#564256',
	    'purple': '#c084fc',
	    'smoky': '#0f1108',
	    'lavender': '#c084fc',
	    'engviolet': '#49416d',
	    'coral': '#e08d79',
	    'napyellow':'#efcb68',
    },
    fontFamily: {
	    sans: ['Nunito Sans', 'sans-serif'],
    },
    extend: {},
  },
  plugins: [],
}

