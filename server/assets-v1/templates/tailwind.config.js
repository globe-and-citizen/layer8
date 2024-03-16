/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{html,js}", "./public/welcome.html"],
  theme: {
    extend: {
      colors: {
        primary: '#ff0000',
        secondary: '#00ff00',
      },
      fontFamily: {
        sans: ['Roboto', 'Arial', 'sans-serif'],
      },
      spacing: {
        '72': '18rem',
        '84': '21rem',
      },
      fontSize: {
        mammoth: '8rem'
      },
      animation: {
        wiggle: 'wiggle 1s ease-in-out infinite',
        animation: 'bounce 2s infinite',
        'slideFromRight': 'slideFromRight 0.5s ease-in-out',
        'slideFromLeft': 'slideFromLeft 0.5s ease-in-out',
      },
      keyframes: {
        wiggle: {
          "0% 100%": {
            transform: "rotate(-1deg)",
          },
          "50%": { transform: "rotate(1deg)" },
        },
        bounce: {
          "0%, 100%": {
            transform: "translateY(-5%)",
            // animation-timing-function: "cubic-bezier(0.8, 0, 1, 1)",
          },
          "50%": {
            transform: "translateY(0)",
            // animation-timing-function: "cubic-bezier(0, 0, 0.2, 1)",
          }
        },
        slideFromRight: {
          '0%': { transform: 'translateX(100%)', opacity: '0' },
          '100%': { transform: 'translateX(0)', opacity: '1' },
        },
        slideFromLeft: {
          '0%': { transform: 'translateX(-100%)', opacity: '0' },
          '100%': { transform: 'translateX(0)', opacity: '1' },
        },
      }
    },
  },
  plugins: [],
}