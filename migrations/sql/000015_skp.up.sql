-- Tabel kategori besar
CREATE TABLE IF NOT EXISTS skp_kategori (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    nama TEXT NOT NULL
);

-- Tabel subkategori (punya foreign key ke kategori)
CREATE TABLE IF NOT EXISTS skp_subkategori (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    kategori_id UUID NOT NULL REFERENCES skp_kategori(id) ON DELETE CASCADE,
    nama TEXT NOT NULL
);

-- Tabel intervensi/item (punya foreign key ke subkategori)
CREATE TABLE IF NOT EXISTS skp_intervensi (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    kategori_id UUID NOT NULL REFERENCES skp_kategori(id) ON DELETE CASCADE,
    subkategori_id UUID NOT NULL REFERENCES skp_subkategori(id) ON DELETE CASCADE,
    nama TEXT NOT NULL
);