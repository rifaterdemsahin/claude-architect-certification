import os
import sys

_kv_client = None

def _get_keyvault_client():
    global _kv_client
    if _kv_client is not None:
        return _kv_client

    vault_name = os.environ.get("AZURE_KEYVAULT_NAME")
    if not vault_name:
        return None

    try:
        from azure.identity import DefaultAzureCredential
        from azure.keyvault.secrets import SecretClient
        
        vault_url = f"https://{vault_name}.vault.azure.net"
        credential = DefaultAzureCredential()
        _kv_client = SecretClient(vault_url=vault_url, credential=credential)
        return _kv_client
    except ImportError:
        print("[KeyVault] Warning: azure-identity or azure-keyvault-secrets not installed. Falling back to local env.", file=sys.stderr)
        return None
    except Exception as e:
        print(f"[KeyVault] Error initializing client: {e}. Falling back to local env.", file=sys.stderr)
        return None

def get_secret(secret_name: str, fallback_env_name: str = None) -> str:
    """
    Retrieves a secret from Azure Key Vault. If Key Vault is not configured
    or retrieval fails, it falls back to the specified environment variable.
    """
    client = _get_keyvault_client()
    if client:
        try:
            # Key Vault secrets only support alphanumeric characters and dashes.
            # Convert name to dash-cased (e.g. ANTHROPIC_API_KEY -> ANTHROPIC-API-KEY)
            kv_secret_name = secret_name.replace("_", "-")
            secret = client.get_secret(kv_secret_name)
            if secret and secret.value:
                return secret.value
        except Exception as e:
            print(f"[KeyVault] Warning: Failed to retrieve secret '{secret_name}' from Key Vault: {e}", file=sys.stderr)

    # Fallback to local environment variables
    env_name = fallback_env_name or secret_name
    return os.environ.get(env_name)
