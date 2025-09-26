#!/bin/bash
output=$(golangci-lint run --disable-all -E maintidx ./... 2>&1)

# Prüfen, ob die Ausgabe leer ist
if [ -z "$output" ]; then
    echo "Keine Ausgabe von golangci-lint erhalten. Möglicherweise sind alle MI-Werte über dem Schwellenwert oder es gibt einen Fehler."
    echo "Versuche, die .golangci.yml zu konfigurieren mit 'under: 100' für maintidx."
    exit 1
fi

mi_values=$(echo "$output" | grep -o 'Maintainability Index: [0-9][0-9]*\(\.[0-9]\+\)\?' \
  | sed 's/[^0-9\.]*//')

# Prüfen, ob MI-Werte gefunden wurden
if [ -z "$mi_values" ]; then
    echo "Keine Maintainability Index-Werte gefunden. Stelle sicher, dass maintidx korrekt läuft."
    echo "Ausgabe von golangci-lint:"
    echo "$output"
    exit 1
fi

# Durchschnitt berechnen
sum=0
count=0
for value in $mi_values; do
    sum=$(echo "$sum + $value" | bc -l)
    count=$((count + 1))
done

if [ $count -eq 0 ]; then
    echo "Fehler: Keine gültigen MI-Werte gefunden."
    exit 1
fi

average=$(echo "scale=2; $sum / $count" | bc -l)

# Ergebnis ausgeben
echo "Gefundene MI-Werte: $mi_values"
echo "Durchschnittlicher Maintainability Index: $average"
