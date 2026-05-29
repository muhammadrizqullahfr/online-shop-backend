-- +goose Up
-- +goose StatementBegin

-- 1. Seed Categories
INSERT INTO categories (id, name, slug) VALUES
(1, 'Pakaian', 'pakaian'),
(2, 'Aksesoris', 'aksesoris'),
(3, 'Sepatu', 'sepatu')
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('categories', 'id'), coalesce(max(id), 1)) FROM categories;

-- 2. Seed Users (Password: password123)
INSERT INTO users (email, password, full_name, role, provider) VALUES
('admin@mediatech.com', '$2a$10$EMzbVKL6RORrk0X1.Cv2leiP.3awk.CXkBJxDI.x4vHWCcF.JkPiO', 'Administrator Sistem', 'admin', 'local'),
('customer@mediatech.com', '$2a$10$EMzbVKL6RORrk0X1.Cv2leiP.3awk.CXkBJxDI.x4vHWCcF.JkPiO', 'John Doe', 'customer', 'local')
ON CONFLICT (email) DO NOTHING;

SELECT setval(pg_get_serial_sequence('users', 'id'), coalesce(max(id), 1)) FROM users;

-- 3. Seed Products (Optimized with FTS Keywords)
INSERT INTO products (id, name, description, category_id, price, original_price, image, rating, reviews, stock, in_stock, featured, tags) VALUES

-- Pakaian (category_id: 1)
('a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Kemeja Flanel Kasual Erigo',
 'Kemeja flanel lengan panjang yang tebal dan nyaman dipakai untuk aktivitas kasual sehari-hari khas Erigo. Produk baju kemeja ini menggunakan bahan katun premium tebal yang menjaga tubuh tetap hangat namun tetap adem saat digunakan. Tersedia dalam berbagai warna pilihan pakaian trendi dan modern.',
 1, 250000.00, 299000.00, 'https://images.unsplash.com/photo-1596755094514-f87e34085b2c?w=500&auto=format&fit=crop&q=60', 4.6, 90, 80, true, true,
 ARRAY['kemeja', 'flanel', 'kasual', 'lengan panjang', 'katun', 'pakaian', 'pria', 'erigo']),

('b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Kaos Polos Premium Distro',
 'Kaos polos bahan cotton combed 30s yang lembut di kulit dan tidak mudah melar. Produk baju kaos kasual ini sangat cocok untuk pria maupun wanita sebagai pilihan pakaian kasual harian. Memiliki jahitan rapi dan kuat, tersedia dalam berbagai warna netral serta cerah khas distro.',
 1, 89000.00, 120000.00, 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=500&auto=format&fit=crop&q=60', 4.5, 215, 200, true, true,
 ARRAY['kaos', 'polos', 'cotton', 'combed', 'kasual', 'pakaian', 'pria', 'wanita', 'distro']),

('c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Jaket Bomber Varsity Pria',
 'Jaket bomber varsity gaya sporty yang cocok untuk tampilan kasual maupun semi-formal. Produk baju luar atau jaket pria ini menggunakan bahan polyester berkualitas dengan lapisan dalam hangat dan nyaman. Dilengkapi saku samping serta saku dada fungsional untuk menunjang pakaian fashion Anda.',
 1, 380000.00, 450000.00, 'https://images.unsplash.com/photo-1591047139829-d91aecb6caea?w=500&auto=format&fit=crop&q=60', 4.7, 58, 45, true, true,
 ARRAY['jaket', 'bomber', 'varsity', 'sporty', 'pakaian', 'pria', 'kasual']),

('d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'Celana Chino Slim Fit Katun',
 'Celana chino slim fit berbahan katun stretch yang nyaman dan fleksibel untuk aktivitas harian. Potongan celana panjang modern yang elegan ini sangat cocok dipadukan dengan baju kemeja maupun kaos polos. Pilihan celana pria formal and kasual yang tersedia dalam warna khaki, navy, dan hitam.',
 1, 199000.00, 250000.00, 'https://images.unsplash.com/photo-1473966968600-fa801b869a1a?w=500&auto=format&fit=crop&q=60', 4.4, 103, 120, true, false,
 ARRAY['celana', 'chino', 'slim fit', 'katun', 'stretch', 'pakaian', 'pria', 'formal', 'kasual']),

('e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'Dress Midi Motif Bunga Wanita',
 'Dress midi cantik dengan motif bunga yang feminin dan elegan. Jenis pakaian baju dress wanita ini menggunakan bahan rayon berkualitas yang jatuh dan lembut sehingga membuat siluet tubuh terlihat indah. Cocok untuk acara kasual, hangout, maupun koleksi baju modis fashion wanita.',
 1, 175000.00, 220000.00, 'https://images.unsplash.com/photo-1572804013309-59a88b7e92f1?w=500&auto=format&fit=crop&q=60', 4.8, 77, 60, true, true,
 ARRAY['dress', 'midi', 'motif bunga', 'wanita', 'rayon', 'feminin', 'kasual', 'pakaian']),

('f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'Sweater Rajut Tebal Oversize',
 'Sweater rajut oversize berbahan wool blend yang hangat dan nyaman dipakai saat cuaca dingin. Desain pakaian baju rajut modern yang stylish ini sangat cocok untuk dipadukan dengan celana jeans panjang maupun rok wanita. Menjadi pilihan baju luar hangat bernuansa earth tone.',
 1, 320000.00, 389000.00, 'https://images.unsplash.com/photo-1576566588028-4147f3842f27?w=500&auto=format&fit=crop&q=60', 4.6, 44, 35, true, false,
 ARRAY['sweater', 'rajut', 'oversize', 'wool', 'tebal', 'hangat', 'pakaian', 'wanita', 'pria']),

('0a1b2c3d-4e5f-4a6b-7c8d-9e0f1a2b3c4d', 'Hoodie Fleece Premium Unisex',
 'Hoodie berbahan fleece premium yang super lembut dan hangat dengan desain unisex yang cocok sebagai baju pria maupun baju wanita. Produk jaket hoodie pakaian kasual ini dilengkapi kantong depan fungsional dan tali hoodie adjustable yang sangat nyaman untuk bersantai.',
 1, 299000.00, 350000.00, 'https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=500&auto=format&fit=crop&q=60', 4.7, 129, 90, true, true,
 ARRAY['hoodie', 'fleece', 'premium', 'unisex', 'hangat', 'kasual', 'pakaian', 'pria', 'wanita']),

-- Aksesoris (category_id: 2)
('1b2c3d4e-5f6a-4b7c-8d9e-0f1a2b3c4d5e', 'Topi Baseball Bordir Keren',
 'Topi baseball dengan aksen bordir logo keren di bagian depan. Menggunakan bahan twill aksesoris kuat dan tahan lama dengan strap jenis topi kasual belakang yang mudah disesuaikan ukurannya. Sangat cocok melengkapi fashion topi pelindung kepala harian Anda.',
 2, 85000.00, 100000.00, 'https://images.unsplash.com/photo-1588850561407-ed78c282e89b?w=500&auto=format&fit=crop&q=60', 4.5, 55, 100, true, false,
 ARRAY['topi', 'baseball', 'bordir', 'aksesoris', 'kasual', 'unisex']),

('2c3d4e5f-6a7b-4c8d-9e0f-1a2b3c4d5e6f', 'Tote Bag Kanvas Aesthetic',
 'Tote bag kanvas tebal dengan desain tas aesthetic yang kekinian dan modis. Produk aksesoris tas kanvas ini mampu menampung banyak barang bawaan kuliah atau laptop dengan aman. Dilengkapi kompartemen aksesoris kantong dalam untuk menyimpan barang penting seperti dompet.',
 2, 125000.00, 150000.00, 'https://images.unsplash.com/photo-1544816155-12df9643f363?w=500&auto=format&fit=crop&q=60', 4.7, 88, 75, true, true,
 ARRAY['tote bag', 'kanvas', 'aesthetic', 'aksesoris', 'tas', 'wanita', 'kuliah']),

('3d4e5f6a-7b8c-4d9e-0f1a-2b3c4d5e6f7a', 'Dompet Kulit Asli Minimalis',
 'Dompet kulit asli dengan desain dompet minimalis yang sangat elegan. Produk aksesoris dompet pria dan wanita ini dilengkapi banyak slot penyimpanan kartu dan kompartemen uang kertas lebar. Jahitan dompet rapi menjamin ketahanan produk premium.',
 2, 275000.00, 320000.00, 'https://images.unsplash.com/photo-1548036328-c9fa89d128fa?w=500&auto=format&fit=crop&q=60', 4.9, 41, 30, true, false,
 ARRAY['dompet', 'kulit asli', 'minimalis', 'aksesoris', 'pria', 'wanita', 'premium']),

('4e5f6a7b-8c9d-4e0f-1a2b-3c4d5e6f7a8b', 'Ikat Pinggang Kulit Formal',
 'Ikat pinggang berbahan kulit berkualitas tinggi dengan kepala gesper logam kokoh. Pilihan aksesoris sabuk atau ikat pinggang formal pria yang timeless, sangat cocok melengkapi lingkar celana chino formal maupun celana denim kasual Anda.',
 2, 150000.00, 190000.00, 'https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=500&auto=format&fit=crop&q=60', 4.3, 27, 50, true, false,
 ARRAY['ikat pinggang', 'belt', 'kulit', 'formal', 'aksesoris', 'pria']),

-- Sepatu (category_id: 3)
('5f6a7b8c-9d0e-4f1a-2b3c-4d5e6f7a8b9c', 'Sneakers Kasual Putih Bersih',
 'Sneakers kasual berwarna putih bersih dengan desain sepatu timeless yang tidak pernah ketinggalan zaman. Menggunakan sol karet sepatu anti-slip empuk yang memberikan kenyamanan melangkah sepanjang hari. Model sepatu sneakers unisex pria dan wanita.',
 3, 450000.00, 520000.00, 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=500&auto=format&fit=crop&q=60', 4.8, 143, 40, true, true,
 ARRAY['sneakers', 'kasual', 'putih', 'sepatu', 'pria', 'wanita', 'anti-slip']),

('6a7b8c9d-0e1f-4a2b-3c4d-5e6f7a8b9c0d', 'Sandal Slide Empuk Kasual',
 'Sandal slide dengan bahan materi EVA yang sangat empuk dan ringan di kaki. Desain alas kaki sandal minimalis kasual yang cocok dipakai santai di rumah maupun aktivitas luar ruangan. Sol sandal anti-slip aman dipadukan dengan celana pendek.',
 3, 99000.00, 130000.00, 'https://images.unsplash.com/photo-1603487742131-4160ec999306?w=500&auto=format&fit=crop&q=60', 4.4, 68, 80, true, false,
 ARRAY['sandal', 'slide', 'empuk', 'kasual', 'sepatu', 'pria', 'wanita', 'eva']),

('7b8c9d0e-1f2a-4b3c-4d5e-6f7a8b9c0d1e', 'Sepatu Boots Kulit Pria',
 'Sepatu boots berbahan kulit asli tangguh yang menampilkan impresi maskulin. Desain produk sepatu klasik ini sangat serasi disandingkan dengan pakaian kemeja formal maupun setelan celana panjang. Sol karet sepatu tebal nyaman seharian.',
 3, 650000.00, 750000.00, 'https://images.unsplash.com/photo-1638247025967-b4e38f787b76?w=500&auto=format&fit=crop&q=60', 4.6, 32, 20, true, true,
 ARRAY['boots', 'kulit', 'pria', 'sepatu', 'maskulin', 'formal', 'kasual', 'premium'])

ON CONFLICT (id) DO NOTHING;

-- 4. Seed Product Variants
INSERT INTO product_variants (id, product_id, name, type, value, price, price_adjustment, stock) VALUES

-- Kemeja Flanel Kasual Erigo
('8c9d0e1f-2a3b-4c4d-5e6f-7a8b9c0d1e2f', 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Ukuran S', 'size', 'S', 250000.00, 0.00, 15),
('9d0e1f2a-3b4c-4d5e-6f7a-8b9c0d1e2f3a', 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Ukuran M', 'size', 'M', 250000.00, 0.00, 25),
('0e1f2a3b-4c5d-4e6f-7a8b-9c0d1e2f3a4b', 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Ukuran L', 'size', 'L', 250000.00, 0.00, 20),
('1f2a3b4c-5d6e-4f7a-8b9c-0d1e2f3a4b5c', 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Ukuran XL', 'size', 'XL', 250000.00, 0.00, 15),
('2a3b4c5d-6e7f-4a8b-9c0d-1e2f3a4b5c6d', 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Ukuran XXL', 'size', 'XXL', 260000.00, 10000.00, 5),

-- Kaos Polos Premium Distro
('3b4c5d6e-7f8a-4b9c-0d1e-2f3a4b5c6d7e', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Ukuran S', 'size', 'S', 89000.00, 0.00, 40),
('4c5d6e7f-8a9b-4c0d-e1f2-3a4b5c6d7e8f', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Ukuran M', 'size', 'M', 89000.00, 0.00, 60),
('5d6e7f8a-9b0c-4d1e-2f3a-4b5c6d7e8f9a', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Ukuran L', 'size', 'L', 89000.00, 0.00, 55),
('6e7f8a9b-0c1d-4e2f-3a4b-5c6d7e8f9a0b', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Ukuran XL', 'size', 'XL', 89000.00, 0.00, 30),
('7f8a9b0c-1d2e-4f3a-4b5c-6d7e8f9a0b1c', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Ukuran XXL', 'size', 'XXL', 99000.00, 10000.00, 15),

-- Jaket Bomber Varsity Pria
('8a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2d', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Ukuran S', 'size', 'S', 380000.00, 0.00, 8),
('9a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2e', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Ukuran M', 'size', 'M', 380000.00, 0.00, 15),
('aa9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2f', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Ukuran L', 'size', 'L', 380000.00, 0.00, 12),
('ba9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2a', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Ukuran XL', 'size', 'XL', 380000.00, 0.00, 7),
('ca9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2b', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Ukuran XXL', 'size', 'XXL', 395000.00, 15000.00, 3),

-- Celana Chino Slim Fit
('da9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2c', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'Ukuran S', 'size', 'S', 199000.00, 0.00, 20),
('ea9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2d', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'Ukuran M', 'size', 'M', 199000.00, 0.00, 35),
('fa9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2e', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'Ukuran L', 'size', 'L', 199000.00, 0.00, 30),
('0a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2f', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'Ukuran XL', 'size', 'XL', 199000.00, 0.00, 25),
('1a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2a', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'Ukuran XXL', 'size', 'XXL', 209000.00, 10000.00, 10),

-- Dress Midi Motif Bunga
('2a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2b', 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'Ukuran S', 'size', 'S', 175000.00, 0.00, 15),
('3a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2c', 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'Ukuran M', 'size', 'M', 175000.00, 0.00, 20),
('4a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2d', 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'Ukuran L', 'size', 'L', 175000.00, 0.00, 15),
('5a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2e', 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'Ukuran XL', 'size', 'XL', 175000.00, 0.00, 10),

-- Sweater Rajut Tebal Oversize
('6a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2f', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'Ukuran S', 'size', 'S', 320000.00, 0.00, 8),
('7a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2a', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'Ukuran M', 'size', 'M', 320000.00, 0.00, 12),
('8a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2b', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'Ukuran L', 'size', 'L', 320000.00, 0.00, 10),
('9a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2c', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'Ukuran XL', 'size', 'XL', 320000.00, 0.00, 5),

-- Hoodie Fleece Premium Unisex
('0a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2d', '0a1b2c3d-4e5f-4a6b-7c8d-9e0f1a2b3c4d', 'Ukuran S', 'size', 'S', 299000.00, 0.00, 20),
('1a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2e', '0a1b2c3d-4e5f-4a6b-7c8d-9e0f1a2b3c4d', 'Ukuran M', 'size', 'M', 299000.00, 0.00, 30),
('2a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2f', '0a1b2c3d-4e5f-4a6b-7c8d-9e0f1a2b3c4d', 'Ukuran L', 'size', 'L', 299000.00, 0.00, 25),
('3a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2a', '0a1b2c3d-4e5f-4a6b-7c8d-9e0f1a2b3c4d', 'Ukuran XL', 'size', 'XL', 299000.00, 0.00, 10),
('4a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2b', '0a1b2c3d-4e5f-4a6b-7c8d-9e0f1a2b3c4d', 'Ukuran XXL', 'size', 'XXL', 315000.00, 16000.00, 5),

-- Sneakers Kasual Putih
('5a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2c', '5f6a7b8c-9d0e-4f1a-2b3c-4d5e6f7a8b9c', 'Ukuran 38', 'size', '38', 450000.00, 0.00, 8),
('6a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2d', '5f6a7b8c-9d0e-4f1a-2b3c-4d5e6f7a8b9c', 'Ukuran 39', 'size', '39', 450000.00, 0.00, 10),
('7a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2e', '5f6a7b8c-9d0e-4f1a-2b3c-4d5e6f7a8b9c', 'Ukuran 40', 'size', '40', 450000.00, 0.00, 10),
('8a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2f', '5f6a7b8c-9d0e-4f1a-2b3c-4d5e6f7a8b9c', 'Ukuran 41', 'size', '41', 450000.00, 0.00, 8),
('9a9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2a', '5f6a7b8c-9d0e-4f1a-2b3c-4d5e6f7a8b9c', 'Ukuran 42', 'size', '42', 450000.00, 0.00, 4),

-- Sandal Slide Empuk
('0b9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2d', '6a7b8c9d-0e1f-4a2b-3c4d-5e6f7a8b9c0d', 'Ukuran 37', 'size', '37', 99000.00, 0.00, 20),
('1b9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2e', '6a7b8c9d-0e1f-4a2b-3c4d-5e6f7a8b9c0d', 'Ukuran 38', 'size', '38', 99000.00, 0.00, 20),
('2b9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2f', '6a7b8c9d-0e1f-4a2b-3c4d-5e6f7a8b9c0d', 'Ukuran 39', 'size', '39', 99000.00, 0.00, 15),
('3b9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2a', '6a7b8c9d-0e1f-4a2b-3c4d-5e6f7a8b9c0d', 'Ukuran 40', 'size', '40', 99000.00, 0.00, 15),
('4b9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2b', '6a7b8c9d-0e1f-4a2b-3c4d-5e6f7a8b9c0d', 'Ukuran 41', 'size', '41', 99000.00, 0.00, 10),

-- Sepatu Boots Kulit Pria
('5b9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2c', '7b8c9d0e-1f2a-4b3c-4d5e-6f7a8b9c0d1e', 'Ukuran 40', 'size', '40', 650000.00, 0.00, 5),
('6b9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2d', '7b8c9d0e-1f2a-4b3c-4d5e-6f7a8b9c0d1e', 'Ukuran 41', 'size', '41', 650000.00, 0.00, 8),
('7b9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2e', '7b8c9d0e-1f2a-4b3c-4d5e-6f7a8b9c0d1e', 'Ukuran 42', 'size', '42', 650000.00, 0.00, 5),
('8b9b0c1d-2e3f-4a4b-5c6d-7e8f9a0b1c2f', '7b8c9d0e-1f2a-4b3c-4d5e-6f7a8b9c0d1e', 'Ukuran 43', 'size', '43', 650000.00, 0.00, 2)

ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM product_variants;
DELETE FROM products;
DELETE FROM users WHERE email IN ('admin@mediatech.com', 'customer@mediatech.com');
DELETE FROM categories;
-- +goose StatementEnd