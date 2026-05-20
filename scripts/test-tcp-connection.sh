#!/usr/bin/env bash
set -Eeuo pipefail

HOST="${HOST:-localhost}"
PORT="${PORT:-5696}"
TIMEOUT="${TIMEOUT:-3}"

usage() {
  cat <<'EOF'
Usage: scripts/test-tcp-connection.sh [--host HOST] [--port PORT] [--timeout SECONDS]

Teste le transport TCP KMIP demo.

Exemples:
  scripts/test-tcp-connection.sh
  scripts/test-tcp-connection.sh --host 127.0.0.1 --port 5696

Variables d'environnement:
  HOST      Hote TCP cible, defaut: localhost
  PORT      Port TCP cible, defaut: 5696
  TIMEOUT   Timeout en secondes, defaut: 3
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --host)
      HOST="${2:?--host demande une valeur}"
      shift 2
      ;;
    --port)
      PORT="${2:?--port demande une valeur}"
      shift 2
      ;;
    --timeout)
      TIMEOUT="${2:?--timeout demande une valeur}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Argument inconnu: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

if ! command -v nc >/dev/null 2>&1; then
  echo "Erreur: la commande 'nc' est requise pour ce test." >&2
  echo "Installe netcat, puis relance le script." >&2
  exit 127
fi

if ! command -v timeout >/dev/null 2>&1; then
  echo "Erreur: la commande 'timeout' est requise pour eviter un blocage du test." >&2
  exit 127
fi

target="${HOST}:${PORT}"

echo "Test ouverture TCP vers ${target}..."
if ! timeout "${TIMEOUT}" nc -z -w "${TIMEOUT}" "${HOST}" "${PORT}" >/dev/null 2>&1; then
  echo "Echec: impossible de se connecter a ${target}." >&2
  echo "Verifie que le backend est lance avec: go run ./cmd/api" >&2
  exit 1
fi
echo "OK: port TCP ouvert."

# Requete TTLV minimale:
# - Operation = Create
# - ObjectType = SymmetricKey
payload='\x42\x00\x5c\x05\x00\x00\x00\x04\x00\x00\x00\x01'
payload+='\x42\x00\x57\x05\x00\x00\x00\x04\x00\x00\x00\x02'

send_with_nc() {
  local mode="$1"

  case "$mode" in
    q)
      printf '%b' "$payload" | timeout "${TIMEOUT}" nc -q 1 -w "${TIMEOUT}" "${HOST}" "${PORT}"
      ;;
    n)
      printf '%b' "$payload" | timeout "${TIMEOUT}" nc -N -w "${TIMEOUT}" "${HOST}" "${PORT}"
      ;;
    basic)
      printf '%b' "$payload" | timeout "${TIMEOUT}" nc -w "${TIMEOUT}" "${HOST}" "${PORT}"
      ;;
  esac
}

echo "Envoi d'une requete TTLV Create minimale..."
response=""
if response="$(send_with_nc q 2>/dev/null)"; then
  :
elif response="$(send_with_nc n 2>/dev/null)"; then
  :
else
  response="$(send_with_nc basic 2>/dev/null || true)"
fi

if [[ -z "$response" ]]; then
  echo "Echec: aucune reponse recue depuis ${target}." >&2
  exit 1
fi

echo "Reponse recue:"
echo "$response"

if [[ "$response" == *'"ok":true'* ]]; then
  echo "OK: le transport TCP a accepte la requete TTLV."
else
  echo "Echec: le serveur TCP a repondu, mais la requete TTLV n'a pas reussi." >&2
  exit 1
fi
