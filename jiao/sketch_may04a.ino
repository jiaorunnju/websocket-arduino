#include <NewPing.h>

#include<NewPing.h>
#define TRIGGER_PIN 12
#define ECHO_PIN 11
#define MAX_DISTANCE 400
#define N 7
NewPing sonar(TRIGGER_PIN, ECHO_PIN, MAX_DISTANCE);
void setup() {
  Serial.begin(9600);
}

long index = 0;
int Adata[N];
int Bdata[N];

void loop() {

  delay(300);
  unsigned long uS = sonar.ping(); 

  int dis = uS / US_ROUNDTRIP_CM;

  int avg = avgFilter(dis);
  Serial.println(avg);
  //Serial.println("\t");
  index++;
}

int avgFilter(int dis) {
  Adata[index % N] = dis;
  int sum = 0;
  for (int i = 0; i < N; i++) {
    sum += Adata[i];
  }
  int avg = sum / N;
  return avg;
}
