#LOG_DIR="/var/systemreport"

#REPORTE="$LOG_DIR/proc_sortedbycpu.rpt"

ps -eo pid,comm,%cpu,%mem --sort=-%cpu | head -n 6 | awk '{if (NR==1) print "Process"$0; else print "Process"$0}'
