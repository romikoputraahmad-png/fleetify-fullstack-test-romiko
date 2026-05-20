const API_URL = '/api';
const tableBody = document.getElementById('tableBodyReports');
const userRoleSelect = document.getElementById('userRoleSelect');
const btnCreateReport = document.getElementById('btnCreateReport');
const modalCreate = new bootstrap.Modal(document.getElementById('modalCreateReport'));

function getHeaders() {
    return {
        'Content-Type': 'application/json',
        'X-User-ID': userRoleSelect.value 
    };
}

// Mengatur UI berdasarkan Role
function updateUI() {
    // Jika login sebagai APPROVAL (value 2), sembunyikan tombol Buat Laporan
    if (userRoleSelect.value === "2") {
        btnCreateReport.style.display = "none";
    } else {
        btnCreateReport.style.display = "block";
    }
}

// Buka modal saat tombol diklik
btnCreateReport.addEventListener('click', () => {
    document.getElementById('formCreateReport').reset();
    modalCreate.show();
});

// Fungsi Menambah Baris Part di Form
function addPartRow() {
    const list = document.getElementById('itemList');
    const row = document.createElement('div');
    row.className = 'row mb-2 item-row';
    
    // Kita copy HTML select dari baris pertama karena isinya sama statis (untuk efisiensi)
    row.innerHTML = `
        <div class="col-md-8">
            <select class="form-select item-select" required>
                <option value="1">Oli Mesin 10W-40 (PART)</option>
                <option value="2">Filter Oli (PART)</option>
                <option value="3">Kampas Rem Depan (PART)</option>
                <option value="4">Jasa Ganti Oli (SERVICE)</option>
                <option value="5">Jasa Servis Rem (SERVICE)</option>
            </select>
        </div>
        <div class="col-md-4">
            <input type="number" class="form-control item-qty" placeholder="Qty" value="1" min="1" required>
        </div>
    `;
    list.appendChild(row);
}

// F-01: Kirim data laporan ke Backend
async function submitReport() {
    const payload = {
        vehicle_id: parseInt(document.getElementById('vehicleId').value),
        odometer: parseInt(document.getElementById('odometer').value),
        complaint: document.getElementById('complaint').value,
        initial_photo: document.getElementById('initialPhoto').value,
        items: []
    };

    // Ambil semua data part/jasa yang dipilih
    const itemRows = document.querySelectorAll('.item-row');
    itemRows.forEach(row => {
        payload.items.push({
            item_id: parseInt(row.querySelector('.item-select').value),
            quantity: parseInt(row.querySelector('.item-qty').value)
        });
    });

    try {
        const res = await fetch(`${API_URL}/reports`, {
            method: 'POST',
            headers: getHeaders(),
            body: JSON.stringify(payload)
        });
        
        if (res.ok) {
            alert('Laporan berhasil dibuat!');
            modalCreate.hide();
            fetchReports(); // Refresh tabel
        } else {
            const err = await res.json();
            alert('Gagal: ' + err.error);
        }
    } catch (e) {
        alert('Terjadi kesalahan koneksi');
    }
}

// F-02 & F-03: Fungsi Aksi (Approve & Complete)
async function actionReport(id, actionStr) {
    let url = `${API_URL}/reports/${id}/${actionStr}`;
    let method = 'PATCH';
    let body = null;

    if (actionStr === 'complete') {
        const photo = prompt("Masukkan URL Foto Bukti Pengerjaan (Simulasi):", "https://example.com/bukti.jpg");
        if (!photo) return; // Batal jika kosong
        body = JSON.stringify({ proof_photo: photo });
    } else {
        const confirmApprove = confirm("Yakin ingin menyetujui laporan ini?");
        if (!confirmApprove) return;
    }

    try {
        const res = await fetch(url, {
            method: method,
            headers: getHeaders(),
            body: body
        });
        if (res.ok) {
            alert(`Laporan berhasil di-${actionStr}!`);
            fetchReports();
        } else {
            const err = await res.json();
            alert('Gagal: ' + err.error);
        }
    } catch (e) {
        alert('Terjadi kesalahan koneksi');
    }
}

// F-04: Mengambil dan menampilkan data ke tabel
async function fetchReports() {
    updateUI(); // Set tampilan awal sesuai role

    try {
        const response = await fetch(`${API_URL}/reports`, { headers: getHeaders() });
        const result = await response.json();

        while (tableBody.firstChild) {
            tableBody.removeChild(tableBody.firstChild);
        }

        if (result.data && result.data.length > 0) {
            result.data.forEach(report => {
                const tr = document.createElement('tr');

                // Render teks biasa
                ['id', 'vehicle', 'user', 'complaint'].forEach(key => {
                    const td = document.createElement('td');
                    if (key === 'vehicle') td.textContent = report.vehicle ? report.vehicle.license_plate : '-';
                    else if (key === 'user') td.textContent = report.user ? report.user.username : '-';
                    else td.textContent = report[key];
                    tr.appendChild(td);
                });

                // Render Badge Status
                const tdStatus = document.createElement('td');
                const spanStatus = document.createElement('span');
                spanStatus.textContent = report.status;
                if (report.status === 'PENDING_APPROVAL') spanStatus.className = 'badge bg-warning text-dark';
                else if (report.status === 'APPROVED') spanStatus.className = 'badge bg-primary';
                else if (report.status === 'COMPLETED') spanStatus.className = 'badge bg-success';
                tdStatus.appendChild(spanStatus);
                tr.appendChild(tdStatus);

                // Render Tombol Aksi sesuai Role dan Status
                const tdAction = document.createElement('td');
                const currentRole = userRoleSelect.value; // 1 = SA, 2 = APPROVAL

                if (currentRole === "2" && report.status === 'PENDING_APPROVAL') {
                    // Tombol Approve untuk Approval
                    const btnApprove = document.createElement('button');
                    btnApprove.className = 'btn btn-sm btn-success';
                    btnApprove.textContent = 'Approve';
                    btnApprove.onclick = () => actionReport(report.id, 'approve');
                    tdAction.appendChild(btnApprove);
                } else if (currentRole === "1" && report.status === 'APPROVED') {
                    // Tombol Complete untuk SA
                    const btnComplete = document.createElement('button');
                    btnComplete.className = 'btn btn-sm btn-info text-white';
                    btnComplete.textContent = 'Selesaikan Laporan';
                    btnComplete.onclick = () => actionReport(report.id, 'complete');
                    tdAction.appendChild(btnComplete);
                } else {
                    tdAction.textContent = "-";
                }
                tr.appendChild(tdAction);

                tableBody.appendChild(tr);
            });
        } else {
            const tr = document.createElement('tr');
            const td = document.createElement('td');
            td.colSpan = 6;
            td.className = "text-center text-muted";
            td.textContent = "Belum ada laporan pemeliharaan.";
            tr.appendChild(td);
            tableBody.appendChild(tr);
        }
    } catch (error) {
        console.error("Gagal mengambil data:", error);
    }
}

document.addEventListener("DOMContentLoaded", fetchReports);
userRoleSelect.addEventListener('change', fetchReports);

// Fitur Bonus B-01: Export data ke CSV tanpa library
async function exportToCSV() {
    try {
        // Ambil data terbaru dari API
        const response = await fetch(`${API_URL}/reports`, { headers: getHeaders() });
        const result = await response.json();

        if (!result.data || result.data.length === 0) {
            alert("Tidak ada data untuk diexport!");
            return;
        }

        // 1. Buat Header Kolom CSV
        let csvContent = "ID Laporan,Nomor Polisi,Pembuat (SA),Odometer (KM),Keluhan,Status,Tanggal Dibuat\n";

        // 2. Looping data dan gabungkan dengan koma
        result.data.forEach(report => {
            const id = report.id;
            const nopol = report.vehicle ? report.vehicle.license_plate : "-";
            const sa = report.user ? report.user.username : "-";
            const odometer = report.odometer;
            // Bersihkan keluhan dari koma atau enter agar tidak merusak format CSV
            const keluhan = `"${report.complaint.replace(/"/g, '""').replace(/\n/g, ' ')}"`;
            const status = report.status;
            const tanggal = new Date(report.created_at).toLocaleString('id-ID');

            csvContent += `${id},${nopol},${sa},${odometer},${keluhan},${status},${tanggal}\n`;
        });

        // 3. Buat Blob dan Link untuk mengunduh file
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const url = URL.createObjectURL(blob);
        
        const link = document.createElement("a");
        link.setAttribute("href", url);
        link.setAttribute("download", "Riwayat_Laporan_Fleetify.csv");
        
        // Simulasi klik secara tersembunyi
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);

    } catch (error) {
        console.error("Gagal mengexport CSV:", error);
        alert("Terjadi kesalahan saat mengunduh CSV");
    }
}