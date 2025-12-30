/** @type {import('jest').Config} */
module.exports = {
  testEnvironment: "jsdom",
  testMatch: ["**/*.test.js"],
  setupFilesAfterEnv: ["<rootDir>/src/test/setup.js"],
  moduleNameMapper: {
    "^\\$lib/(.*)$": "<rootDir>/src/lib/$1",
    "^svelte$": "<rootDir>/node_modules/svelte",
  },
  transform: {
    "^.+\\.svelte$": ["svelte-jest", { preprocess: true }],
    "^.+\\.js$": "babel-jest",
  },
  moduleFileExtensions: ["js", "svelte"],
  testTimeout: 10000,
};
