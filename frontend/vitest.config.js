import { defineConfig } from "vitest/config";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  plugins: [
    svelte({
      compilerOptions: {
        runes: true,
      },
    }),
  ],
  resolve: {
    alias: {
      $lib: path.resolve(__dirname, "./src/lib"),
    },
    // Force browser conditions for Svelte 5 client-side rendering
    conditions: ["browser"],
  },
  test: {
    globals: true,
    include: ["src/**/*.test.js"],
    setupFiles: ["src/test/setup.js"],
    environment: "happy-dom",
    // Vitest 4: inline Svelte for proper resolution
    server: {
      deps: {
        inline: [/svelte/, "@testing-library/svelte"],
      },
    },
  },
});
