module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',  // Sesuaikan dengan alamat server Go Anda
        changeOrigin: true,
        pathRewrite: {
          '^/api': ''
        }
      }
    }
  }
}
