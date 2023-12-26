import path from 'path'
import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
    build: {
        copyPublicDir: false,
        emptyOutDir: true,
        modulePreload: false,
        chunkSizeWarningLimit: 1600,
        rollupOptions: {
            input: {
                admin: path.resolve(__dirname, 'html/admin.html'),
                profile: path.resolve(__dirname, 'html/profile.html'),
                error: path.resolve(__dirname, 'html/error.html'),
                signin: path.resolve(__dirname, 'html/auth.html'),
            }
        },
    },
    server: {
        host: '0.0.0.0',
        port: 8080,
        hmr: true,
    },
    plugins: [vue({isProduction: false})],
})
