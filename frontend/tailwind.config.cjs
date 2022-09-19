/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './index.html',
    './src/**/*.{html,js,ts,jsx,tsx}'
  ],
  theme: {
    extend: {
      colors: {
        neutral: {
          // Neutral color: shades of gray, brown, etc.
          // Order from least intense (e.g. for background color) to most intense (e.g. for buttons)
          a: '#000000',
          b: '#333333',
          c: '#666666',
          d: '#999999',
          e: '#cccccc',
          f: '#ffffff'
        },
        main: {
          // Main color: the main color of the brand.
          // Intensity should match with neutral color
          //  if neutral least intense is dark, main least intense should be dark, vice versa.
          // When changing the main color hue, make sure to keep the saturation and lightness the same.
          //  If you're changing the neutral color, the lightness can potentially be changed
          a: '#10451d',
          b: '#1a7431',
          c: '#25a244',
          d: '#2dc653',
          e: '#6ede8a',
          f: '#b7efc5'
        }
      }
    }
  },
  plugins: [] // [require('daisyui')]
}
