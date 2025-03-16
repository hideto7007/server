// DB/sql.go
package DB

const GetIncomeDataInRangeSyntax = `
			SELECT income_forecast_id, payment_date, age, industry, total_amount, deduction_amount, take_home_amount, classification, user_id
			FROM income_forecast_data
			WHERE payment_date BETWEEN $1 AND $2 AND user_id = $3
			ORDER BY payment_date DESC;
			`
const GetDateRangeSyntax = `
			SELECT user_id, MIN(payment_date) as "start_paymaent_date", MAX(payment_date) as "end_paymaent_date" from income_forecast_data
			WHERE user_id = $1
			GROUP BY user_id;
			`
const GetYearsIncomeAndDeductionSyntax = `
			SELECT 
				TO_CHAR(payment_date, 'YYYY') as "year" ,
				SUM(total_amount) as "sum_total_amount", 
				SUM(deduction_amount) as "sum_deduction_amount",  
				SUM(take_home_amount) as "sum_take_home_amount"
			FROM income_forecast_data
			WHERE user_id = $1
			GROUP BY TO_CHAR(payment_date, 'YYYY')
			ORDER BY TO_CHAR(payment_date, 'YYYY') asc;
			`
const InsertIncomeSyntax = `
			INSERT INTO income_forecast_data
			(income_forecast_id, payment_date, age, industry, total_amount, deduction_amount, take_home_amount, created_at, classification, user_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
			`
const UpdateIncomeSyntax = `
			UPDATE income_forecast_data
			SET 
				payment_date = $1, 
				age = $2, 
				industry = $3, 
				total_amount = $4, 
				deduction_amount = $5, 
				take_home_amount = $6, 
				created_at = $7, 
				update_user = $8,
				classification = $9
			WHERE income_forecast_id = $10;
			`

const DeleteIncomeSyntax = `
			DELETE FROM income_forecast_data
			WHERE income_forecast_id = $1;
			`

const GetSignInSyntax = `
			SELECT user_id, user_email, user_password
			FROM users
			WHERE user_email = $1;
			`

const PasswordCheckSyntax = `
			SELECT user_email
			FROM users
			WHERE user_id = $1;
			`

const GetExternalAuthSyntax = `
			SELECT user_id, user_email
			FROM users
			WHERE user_email = $1;
			`

const PostSignUpSyntax = `
			INSERT INTO users
			(user_email, user_password, create_user, create_at, update_user, update_at, delete_flag)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING user_id;
			`

const PutSignInEditSyntax = `
			UPDATE users
			SET
				user_email = coalesce($1, user_email),
				user_password = coalesce($2, user_password),
				update_at  = $3
			WHERE 
				user_id = $4;
			`

const PutPasswordSyntax = `
			UPDATE users
			SET
				user_password = $1,
				update_at  = $2
			WHERE 
				user_id = $3;
			`

const DeleteSignInSyntax = `
			DELETE FROM users
			WHERE user_id = $1
			and user_email = $2;
			`
