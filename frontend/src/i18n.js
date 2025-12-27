import { register, init, getLocaleFromNavigator, locale } from 'svelte-i18n';

const DEFAULT_LOCALE = 'en';

// English translations
const en = {
  app: { name: 'CareKeeper', tagline: 'Healthcare Volunteer Dashboard' },
  auth: {
    login: 'Login', register: 'Register', logout: 'Logout', username: 'Username', password: 'Password',
    confirmPassword: 'Confirm Password', fullName: 'Full Name', loginButton: 'Login', registerButton: 'Create Account',
    loggingIn: 'Logging in...', creatingAccount: 'Creating account...', enterUsername: 'Enter your username',
    enterPassword: 'Enter your password', chooseUsername: 'Choose a username', choosePassword: 'Choose a password',
    confirmYourPassword: 'Confirm your password', yourFullName: 'Your full name', passwordsDoNotMatch: 'Passwords do not match',
    minChars: 'Min {n} characters', connectionError: 'Connection error. Please check your internet connection.',
    loginFailed: 'Login failed', registrationFailed: 'Registration failed. Please try again.'
  },
  common: {
    save: 'Save', cancel: 'Cancel', delete: 'Delete', edit: 'Edit', add: 'Add', close: 'Close',
    loading: 'Loading...', search: 'Search', searchPlaceholder: 'Search...', noResults: 'No results found',
    required: 'Required', optional: 'Optional', yes: 'Yes', no: 'No', confirm: 'Confirm',
    viewAll: 'View all', refresh: 'Refresh', you: 'You', more: 'More', profile: 'Profile', language: 'Language'
  },
  navigation: { dashboard: 'Dashboard', patients: 'Patients', users: 'Users', volunteerDashboard: 'Volunteer Dashboard' },
  dashboard: {
    title: 'Dashboard', totalPatients: 'Total Patients', totalReminders: 'Total Reminders',
    completed: 'Completed', pending: 'Pending', upcomingReminders: 'Upcoming Reminders',
    noUpcomingReminders: 'No upcoming reminders', recentPatients: 'Recent Patients', noPatients: 'No patients yet'
  },
  patients: {
    title: 'Patients', addPatient: 'Add Patient', editPatient: 'Edit Patient', patientName: 'Patient name',
    phone: 'Phone', whatsappNumber: 'WhatsApp Number', whatsappNote: 'WhatsApp number for reminder notifications',
    email: 'Email', emailPlaceholder: 'Email address', notes: 'Notes', notesPlaceholder: 'Any notes about the patient...',
    noPatients: 'No patients found', noPatientsMatch: 'No patients match your search criteria.',
    getStarted: 'Get started by adding your first patient.', noPhone: 'No phone',
    reminders: 'Reminders', addReminder: 'Add reminder', editReminder: 'Edit reminder'
  },
  reminders: {
    title: 'Title', titlePlaceholder: 'e.g., Medication reminder', description: 'Description',
    descriptionPlaceholder: 'Additional details...', dueDate: 'Due Date & Time', priority: 'Priority',
    low: 'Low', medium: 'Medium', high: 'High', recurrence: 'Recurrence', repeat: 'Repeat',
    doesNotRepeat: 'Does not repeat', daily: 'Daily', weekly: 'Weekly', monthly: 'Monthly', yearly: 'Yearly',
    repeatEvery: 'Repeat every', daysOfWeek: 'Days of week', endDate: 'End date (optional)'
  },
  users: {
    title: 'User Management', addUser: 'Add User', editUserRole: 'Edit User Role', user: 'User',
    role: 'Role', created: 'Created', actions: 'Actions', noUsers: 'No users found', noName: 'No name',
    superadmin: 'Superadmin', admin: 'Admin', volunteer: 'Volunteer',
    roleDescription: { admin: 'Admins can view all patients but cannot manage users.', volunteer: 'Volunteers can only see patients they created.' }
  },
  confirm: { deletePatient: 'Are you sure you want to delete this patient?' },
  password: { strength: { weak: 'Weak', fair: 'Fair', good: 'Good', strong: 'Strong' } },
  validation: { usernameMin: 'Username must be at least 3 characters' }
};

// Indonesian translations
const id = {
  app: { name: 'CareKeeper', tagline: 'Dasbor Relawan Kesehatan' },
  auth: {
    login: 'Masuk', register: 'Daftar', logout: 'Keluar', username: 'Nama Pengguna', password: 'Kata Sandi',
    confirmPassword: 'Konfirmasi Kata Sandi', fullName: 'Nama Lengkap', loginButton: 'Masuk', registerButton: 'Buat Akun',
    loggingIn: 'Sedang masuk...', creatingAccount: 'Membuat akun...', enterUsername: 'Masukkan nama pengguna Anda',
    enterPassword: 'Masukkan kata sandi Anda', chooseUsername: 'Pilih nama pengguna', choosePassword: 'Pilih kata sandi',
    confirmYourPassword: 'Konfirmasi kata sandi Anda', yourFullName: 'Nama lengkap Anda', passwordsDoNotMatch: 'Kata sandi tidak cocok',
    minChars: 'Min {n} karakter', connectionError: 'Kesalahan koneksi. Periksa koneksi internet Anda.',
    loginFailed: 'Gagal masuk', registrationFailed: 'Pendaftaran gagal. Silakan coba lagi.'
  },
  common: {
    save: 'Simpan', cancel: 'Batal', delete: 'Hapus', edit: 'Ubah', add: 'Tambah', close: 'Tutup',
    loading: 'Memuat...', search: 'Cari', searchPlaceholder: 'Cari...', noResults: 'Tidak ada hasil',
    required: 'Wajib diisi', optional: 'Opsional', yes: 'Ya', no: 'Tidak', confirm: 'Konfirmasi',
    viewAll: 'Lihat semua', refresh: 'Segarkan', you: 'Anda', more: 'Lainnya', profile: 'Profil', language: 'Bahasa'
  },
  navigation: { dashboard: 'Beranda', patients: 'Pasien', users: 'Pengguna', volunteerDashboard: 'Dasbor Relawan' },
  dashboard: {
    title: 'Beranda', totalPatients: 'Total Pasien', totalReminders: 'Total Pengingat',
    completed: 'Selesai', pending: 'Tertunda', upcomingReminders: 'Pengingat Mendatang',
    noUpcomingReminders: 'Tidak ada pengingat mendatang', recentPatients: 'Pasien Terbaru', noPatients: 'Belum ada pasien'
  },
  patients: {
    title: 'Pasien', addPatient: 'Tambah Pasien', editPatient: 'Ubah Pasien', patientName: 'Nama pasien',
    phone: 'Telepon', whatsappNumber: 'Nomor WhatsApp', whatsappNote: 'Nomor WhatsApp untuk notifikasi pengingat',
    email: 'Email', emailPlaceholder: 'Alamat email', notes: 'Catatan', notesPlaceholder: 'Catatan tentang pasien...',
    noPatients: 'Tidak ada pasien ditemukan', noPatientsMatch: 'Tidak ada pasien yang cocok dengan kriteria pencarian Anda.',
    getStarted: 'Mulai dengan menambahkan pasien pertama Anda.', noPhone: 'Tidak ada telepon',
    reminders: 'Pengingat', addReminder: 'Tambah pengingat', editReminder: 'Ubah pengingat'
  },
  reminders: {
    title: 'Judul', titlePlaceholder: 'mis., Pengingat obat', description: 'Deskripsi',
    descriptionPlaceholder: 'Detail tambahan...', dueDate: 'Tanggal & Waktu Jatuh Tempo', priority: 'Prioritas',
    low: 'Rendah', medium: 'Sedang', high: 'Tinggi', recurrence: 'Pengulangan', repeat: 'Ulangi',
    doesNotRepeat: 'Tidak berulang', daily: 'Harian', weekly: 'Mingguan', monthly: 'Bulanan', yearly: 'Tahunan',
    repeatEvery: 'Ulangi setiap', daysOfWeek: 'Hari dalam seminggu', endDate: 'Tanggal berakhir (opsional)'
  },
  users: {
    title: 'Manajemen Pengguna', addUser: 'Tambah Pengguna', editUserRole: 'Ubah Peran Pengguna', user: 'Pengguna',
    role: 'Peran', created: 'Dibuat', actions: 'Aksi', noUsers: 'Tidak ada pengguna ditemukan', noName: 'Tanpa nama',
    superadmin: 'Superadmin', admin: 'Admin', volunteer: 'Relawan',
    roleDescription: { admin: 'Admin dapat melihat semua pasien tetapi tidak dapat mengelola pengguna.', volunteer: 'Relawan hanya dapat melihat pasien yang mereka buat.' }
  },
  confirm: { deletePatient: 'Apakah Anda yakin ingin menghapus pasien ini?' },
  password: { strength: { weak: 'Lemah', fair: 'Cukup', good: 'Baik', strong: 'Kuat' } },
  validation: { usernameMin: 'Nama pengguna minimal 3 karakter' }
};

register('en', () => Promise.resolve(en));
register('id', () => Promise.resolve(id));

init({
  fallbackLocale: DEFAULT_LOCALE,
  initialLocale: localStorage.getItem('locale') || getLocaleFromNavigator() || DEFAULT_LOCALE
});

export { locale };
