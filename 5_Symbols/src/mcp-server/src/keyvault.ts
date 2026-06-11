import { DefaultAzureCredential } from "@azure/identity";
import { SecretClient } from "@azure/keyvault-secrets";

let kvClient: SecretClient | null = null;

function getKeyvaultClient(): SecretClient | null {
  if (kvClient !== null) {
    return kvClient;
  }

  const vaultName = process.env.AZURE_KEYVAULT_NAME;
  if (!vaultName) {
    return null;
  }

  try {
    const vaultUri = `https://${vaultName}.vault.azure.net`;
    const credential = new DefaultAzureCredential();
    kvClient = new SecretClient(vaultUri, credential);
    return kvClient;
  } catch (error) {
    console.error("[KeyVault] Failed to initialize Azure Key Vault client:", error);
    return null;
  }
}

/**
 * Retrieves a secret from Azure Key Vault. Falls back to local process.env.
 */
export async function getSecret(secretName: string, fallbackEnvName?: string): Promise<string | undefined> {
  const client = getKeyvaultClient();
  if (client) {
    try {
      // Key Vault secret names only support alphanumeric characters and dashes.
      // Convert underscore to dash (e.g., SUPABASE_URL -> SUPABASE-URL)
      const kvSecretName = secretName.replace(/_/g, "-");
      const secret = await client.getSecret(kvSecretName);
      if (secret && secret.value) {
        return secret.value;
      }
    } catch (error) {
      console.warn(`[KeyVault] Warning: Failed to retrieve secret '${secretName}' from Key Vault:`, error);
    }
  }

  // R-05: Fallback to local environment variable; warn if also missing
  const envName = fallbackEnvName || secretName;
  const value = process.env[envName];
  if (value === undefined) {
    console.warn(`[KeyVault] Secret '${secretName}' not found in Key Vault or env. Callers receiving undefined may fail silently.`);
  }
  return value;
}
