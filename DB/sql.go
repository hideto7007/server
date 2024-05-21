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
			INSERT INTO public.income_forecast_data
			(income_forecast_id, payment_date, age, industry, total_amount, deduction_amount, take_home_amount, delete_flag, update_user, created_at, classification, user_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
			`
const UpdateIncomeSyntax = `
			UPDATE public.income_forecast_data
			SET 
				payment_date = $1, 
				age = $2, 
				industry = $3, 
				total_amount = $4, 
				deduction_amount = $5, 
				take_home_amount = $6, 
				created_at = $7, 
				classification = $8
			WHERE income_forecast_id = $9;
			`
const DeleteIncomeSyntax = `
			DELETE FROM public.income_forecast_data
			WHERE income_forecast_id = $1;
			`
