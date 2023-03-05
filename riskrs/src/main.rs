use rand::Rng;
use std::collections::HashMap;

// base cases - 1 defender
// percent attacker wins in 1 vs. 1
const PCT1V1: f64 = 15.0 / 36.0;
// percent attacker wins in 2 vs. 1
const PCT2V1: f64 = 125.0 / 216.0;
// percent attacker wins in 3 vs. 1
const PCT3V1: f64 = 855.0 / 1296.0;

// base cases - 2 defenders
// percent attacker wins in 1 vs. 2
const PCT1V2: f64 = 55.0 / 216.0;
// percent attacker wins 2 in 2 vs. 2
const PCT2V2_AW2: f64 = 295.0 / 1296.0;
// percent attacker splits 1:1 with defender in 2 vs. 2
const PCT2V2_SP: f64 = 420.0 / 1296.0;
// percent attacker wins 2 in 3 vs. 2
const PCT3V2_AW2: f64 = 2890.0 / 7776.0;
// percent attacker loses 2 in 3 vs. 2
const PCT3V2_DW2: f64 = 2275.0 / 7776.0;
// percent attacker splits 1:1 with defender in 3 vs. 2
const PCT3V2_SP: f64 = 2611.0 / 7776.0;

// battle_prob returns the probability that attacker wins
fn battle_prob(a: u8, d: u8, cache: &mut HashMap<(u8, u8), f64>) -> f64 {
    if let Some(v) = cache.get(&(a, d)) {
        return *v;
    }

    if a == 0 && d == 0 {
        panic!("unreachable: a & d must not both be 0!");
    }

    // declare as anonymous fn so can cache the result easily
    let ret = || -> f64 {
        if a == 0 {
            return 0.0;
        }
        if d == 0 {
            return 1.0;
        }

        // only 1 defender remains - battle may end
        if d == 1 {
            if a == 1 {
                return PCT1V1;
            }
            if a == 2 {
                return PCT2V1 + (1.0 - PCT2V1) * battle_prob(a - 1, d, cache);
            }
            return PCT3V1 + (1.0 - PCT3V1) * battle_prob(a - 1, d, cache);
        }

        // only 2 defenders remain - battle may end
        if d == 2 && a < 3 {
            if a == 1 {
                // 1 vs. 2 requires attacker to win in 2 turns
                return PCT1V2 * battle_prob(a, d - 1, cache);
            }
            if a == 2 {
                return PCT2V2_AW2 + PCT2V2_SP * battle_prob(a - 1, d - 1, cache);
            }
        }

        // otherwise - many defenders remain
        if a == 1 {
            return PCT1V2 * battle_prob(a, d - 1, cache);
        }
        if a == 2 {
            return PCT2V2_AW2 * battle_prob(a, d - 2, cache)
                + PCT2V2_SP * battle_prob(a - 1, d - 1, cache);
        }
        return PCT3V2_AW2 * battle_prob(a, d - 2, cache)
            + PCT3V2_DW2 * battle_prob(a - 2, d, cache)
            + PCT3V2_SP * battle_prob(a - 1, d - 1, cache);
    }();

    cache.insert((a, d), ret);
    return ret;
}

// roll simulates an attack (i.e. a roll of the dice). it accepts the number of dice / troops and returns the number remaining after simulation.
fn roll(attackers: u8, defenders: u8) -> (u8, u8) {
    if attackers == 0 || attackers > 3 {
        panic!("must have 1-3 attackers");
    }
    if defenders == 0 || defenders > 2 {
        panic!("must have 1 or 2 defenders");
    }

    let count = min(attackers, defenders);
    if count != 1 && count != 2 {
        panic!("expected battle of 1 or 2");
    }

    let mut rng = rand::thread_rng();
    let mut rolls_a: Vec<u8> = (0..attackers).map(|_| rng.gen_range(0..7)).collect();
    let mut rolls_b: Vec<u8> = (0..defenders).map(|_| rng.gen_range(0..7)).collect();

    rolls_a.sort_by(|a, b| a.cmp(b).reverse());
    rolls_b.sort_by(|a, b| a.cmp(b).reverse());

    let mut remaining_a = attackers;
    let mut remaining_b = defenders;

    for i in 0..count {
        let a = rolls_a.get(usize::from(i)).unwrap();
        let d = rolls_b.get(usize::from(i)).unwrap();

        if a > d {
            remaining_b -= 1;
        } else {
            remaining_a -= 1;
        }
    }

    (remaining_a, remaining_b)
}

fn min(a: u8, b: u8) -> u8 {
    if a < b {
        return a;
    }
    return b;
}

fn clamp(n: u8, min: u8, max: u8) -> u8 {
    if n < min {
        return min;
    } else if n > max {
        return max;
    }
    n
}

fn simulate_battle(attackers: u8, defenders: u8) -> bool {
    let mut remaining_a = attackers;
    let mut remaining_d = defenders;

    while remaining_a > 0 && remaining_d > 0 {
        let send_a = clamp(remaining_a, 1, 3);
        let send_d = clamp(remaining_d, 1, 2);

        let (ret_a, ret_d) = roll(send_a, send_d);
        remaining_a -= send_a - ret_a;
        remaining_d -= send_d - ret_d;
    }
    return remaining_a > 0;
}

fn main() {
    println!("battle probabilties ---------------------------------------------------------------");
    let size = 35u8;
    let mut cache: HashMap<(u8, u8), f64> = HashMap::new();
    for d in 1..(size + 1) {
        for a in 1..(size + 1) {
            print!("{:>7.3}%\t", battle_prob(a, d, &mut cache) * 100.0)
        }
        println!("");
    }

    let size_sim = 20u8;
    let trials = 1000u32;
    println!("battle probabilties (simulated) ---------------------------------------------------------------");
    for d in 1..size_sim + 1 {
        for a in 1..size_sim + 1 {
            let mut a_wins = 0;
            for _ in 0..trials {
                if simulate_battle(a, d) {
                    a_wins += 1;
                }
            }
            print!("{:>7.3}%\t", f64::from(a_wins) / f64::from(trials) * 100.0);
        }
        println!("");
    }
    println!("");
}
