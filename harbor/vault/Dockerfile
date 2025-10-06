# Gunakan base image resmi Vault
FROM hashicorp/vault:1.16.2

# Set environment variable supaya Vault jalan di server mode
ENV VAULT_ADDR=http://0.0.0.0:8200
ENV VAULT_API_ADDR=http://0.0.0.0:8200

# Buat folder config di dalam container
RUN mkdir -p /vault/config

# Copy konfigurasi Vault ke container
COPY ./vault.hcl /vault/config/vault.hcl

# Expose port Vault
EXPOSE 8200

# Jalankan Vault dalam mode server dengan config yang kita buat
ENTRYPOINT ["vault", "server", "-config=/vault/config/vault.hcl"]
