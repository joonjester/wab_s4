import matplotlib.pyplot as plt
import json
from collections import Counter

results = "results.json"

with open(results, "r") as f:
    data = json.load(f)

features = list(data.keys())

# --- Test Coverage ---
test_coverage = [float(data[f]["test_coverage"].replace("%","")) for f in features]
plt.figure(figsize=(8,5))
plt.bar(features, test_coverage, color='skyblue')
plt.ylabel("Test Coverage (%)")
plt.title("Test Coverage pro Feature")
plt.grid(axis="y", linestyle="--", alpha=0.7)
plt.show()

# --- Mutation Score ---
mutation_score = [float(data[f]["mutation_score"].replace("%","")) for f in features]
plt.figure(figsize=(8,5))
plt.bar(features, mutation_score, color='lightgreen')
plt.ylabel("Mutation Score (%)")
plt.title("Mutation Score pro Feature")
plt.grid(axis="y", linestyle="--", alpha=0.7)
plt.show()

# --- Benchmark Zeit (Nanosecond pro Operation) ---
benchmark_sec = [float(data[f]["benchmark_score"].replace("ns/op","")) for f in features]
plt.figure(figsize=(8,5))
plt.bar(features, benchmark_sec, color='salmon')
plt.ylabel("Benchmark Zeit (ns/op)")
plt.title("Benchmark Zeit pro Feature")
plt.grid(axis="y", linestyle="--", alpha=0.7)
plt.show()

# --- Cyclo ---
plt.figure(figsize=(10, 6))
for feature, metrics in data.items():
    cyclo = metrics["cyclo_score"]
    counter = Counter(cyclo)
    
    x_vals = list(counter.keys())
    y_vals = [int(feature[-1])*10 for _ in x_vals]  
    sizes = [counter[k]*30 for k in x_vals]        
    plt.scatter(x_vals, y_vals, s=sizes, alpha=0.6, label=feature)

plt.yticks([10, 20, 30], ["feature1", "feature2", "feature3"])
plt.xlabel("Cyclomatic Complexity")
plt.ylabel("Feature")
plt.title("Cyclomatic Complexity Scatter Plot (Point Size ~ Frequency)")
plt.legend()
plt.grid(True)
plt.show()

# --- MI ---
plt.figure(figsize=(10, 6))
for feature, metrics in data.items():
    mi = metrics["maintainability_index"]
    counter = Counter(mi)
    
    # Scatter-Plot mit Punktgröße proportional zur Häufigkeit
    x_vals = list(counter.keys())
    y_vals = [int(feature[-1])*10 for _ in x_vals]  # z.B. Feature als y-Achse (Feature1->10, Feature2->20)
    sizes = [counter[k]*30 for k in x_vals]          # Punktgröße = Anzahl * 30
    plt.scatter(x_vals, y_vals, s=sizes, alpha=0.6, label=feature)

plt.yticks([10, 20, 30], ["feature1", "feature2", "feature3"])
plt.xlabel("Maintainability Index")
plt.ylabel("Feature")
plt.title("Maintainability Index Scatter Plot (Point Size ~ Frequency)")
plt.legend()
plt.grid(True)
plt.show()
