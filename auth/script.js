/**
 * Logika Autentikasi Petelot Coding
 * Mengatur transisi instan antara Login dan Register
 */
function switchMode() {
    const container = document.getElementById('mainContainer');
    const loginArea = document.getElementById('loginArea');
    const registerArea = document.getElementById('registerArea');

    // Memulai transisi sliding dengan menukar class pada container
    if (container.classList.contains('login-mode')) {
        container.classList.remove('login-mode');
        container.classList.add('register-mode');
        
        // Pergantian teks instan: sembunyikan login, tampilkan register
        loginArea.style.display = 'none';
        registerArea.style.display = 'block';
        registerArea.style.opacity = '1'; 
    } else {
        container.classList.remove('register-mode');
        container.classList.add('login-mode');
        
        // Pergantian teks instan: sembunyikan register, tampilkan login
        registerArea.style.display = 'none';
        loginArea.style.display = 'block';
        loginArea.style.opacity = '1';
    }
}

/**
 * Logika Matrix Rain Background
 * Berjalan di satu layar penuh secara stabil
 */
const canvas = document.getElementById('matrixCanvas');
const ctx = canvas.getContext('2d');
let columns, drops = [];
const fontSize = 16;
const characters = '01ABCDEFGHIJKLMNOPQRSTUVWXYZ$#{}[]';

// Inisialisasi ukuran canvas sesuai lebar jendela browser
function initCanvas() {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    columns = Math.floor(canvas.width / fontSize);
    drops = Array(columns).fill(1);
}

// Fungsi menggambar tetesan Matrix
function draw() {
    // Memberikan efek trail hitam transparan
    ctx.fillStyle = 'rgba(0, 0, 0, 0.08)';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Warna teks rain putih agar kontras dengan background hitam
    ctx.fillStyle = '#fff';
    ctx.font = fontSize + 'px monospace';

    for (let i = 0; i < drops.length; i++) {
        const text = characters.charAt(Math.floor(Math.random() * characters.length));
        ctx.fillText(text, i * fontSize, drops[i] * fontSize);

        // Reset tetesan ke atas setelah mencapai batas bawah layar
        if (drops[i] * fontSize > canvas.height && Math.random() > 0.975) {
            drops[i] = 0;
        }
        drops[i]++;
    }
}

// Listener untuk mengatur ulang canvas saat ukuran jendela berubah
window.addEventListener('resize', initCanvas);

// Memulai animasi Matrix Rain
initCanvas();
setInterval(draw, 35);

// Menghindari refresh halaman pada pengiriman form untuk keperluan pengujian
const loginForm = document.getElementById('loginForm');
if (loginForm) {
    loginForm.onsubmit = (e) => { 
        e.preventDefault(); 
        console.log('Login attempt recorded.'); 
    };
}

const registerForm = document.getElementById('registerForm');
if (registerForm) {
    registerForm.onsubmit = (e) => { 
        e.preventDefault(); 
        console.log('Registration attempt recorded.'); 
    };
}
