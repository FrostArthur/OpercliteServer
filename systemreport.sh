# Carpeta donde se guardarán los reportes
# Agregar cantidad de procesos en ejecucion
LOG_DIR="/var/systemreport"

#Archivo con marca de tiempo
REPORTE="$LOG_DIR/systemreport.rpt"

# === CPU ===
# Obtenemos el porcentaje de CPU usado (us + sy)
cpu_usage=$(top -bn1 | grep "Cpu(s)" | awk '{print 100 - $8}')
cpu_usage=$(printf "%.1f" "$cpu_usage")

# === MEMORIA ===
# Tomamos memoria total y usada de free -m (en MiB)
read mem_total mem_used <<< $(free -m | awk '/Mem:/ {print $2, $3}')
mem_percent=$(awk "BEGIN {printf \"%.1f\", ($mem_used/$mem_total)*100}")
mem_used_gb=$(awk "BEGIN {printf \"%.1f\", $mem_used/1024}")
mem_total_gb=$(awk "BEGIN {printf \"%.1f\", $mem_total/1024}")

# === DISCO ===
# Calculamos promedio de uso de disco raíz (/) — o de todos si prefieres
disk_usage=$(df -h --output=pcent / | tail -1 | tr -dc '0-9')

# === CANTIDAD DE PROCESOS EN EJECUCION ====
number_process=$(ps -aux | wc -l)

# === FORMATO DE SALIDA ===
echo "CPU: ${cpu_usage} | MEM: ${mem_percent} | (${mem_used_gb} GiB/${mem_total_gb} GiB) | DISK: ${disk_usage} | ${number_process} PROCESSES" > "$REPORTE"
