-- name: ListKategoriSubkategoriIntervensi :many
SELECT
    k.id AS kategori_id,
    k.nama AS kategori_nama,
    s.id AS subkategori_id,
    s.nama AS subkategori_nama,
    i.id AS intervensi_id,
    i.nama AS intervensi_nama
FROM
    public.skp_intervensi i
JOIN
    public.skp_subkategori s ON i.subkategori_id = s.id
JOIN
    public.skp_kategori k ON i.kategori_id = k.id
ORDER BY
    k.nama, s.nama, i.nama;